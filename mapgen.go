package main

func (call *IrInterfaceMethodCall) emit(args []Expr) {
	emit("# emit interface method call \"%s\"", call.methodName)
	mapType := &Gtype{
		kind: G_MAP,
		mapKey: &Gtype{
			kind: G_STRING,
		},
		mapValue: &Gtype{
			kind: G_STRING,
		},
	}
	emit("# emit receiverTypeId of %s", call.receiver.getGtype().String())
	emitOffsetLoad(call.receiver, ptrSize, ptrSize)
	emit("IMUL_NUMBER 8")
	emit("PUSH_PRIMITIVE")

	emit("lea receiverTypes(%%rip), %%rax")
	emit("PUSH_PRIMITIVE")
	emit("SUM_FROM_STACK")

	emit("# find method %s", call.methodName)
	emit("mov (%%rax), %%r10") // address of receiverType

	emit("mov $128, %%r11")  // copy len

	emit("lea .M%s, %%rax", call.methodName) // index value
	emit("mov %%rax, %%r12")                 // index value
	emitMapGet(mapType, false)

	emit("PUSH_PRIMITIVE")

	emit("# setting arguments (len=%d)", len(args))

	receiver := args[0]
	emit("mov $0, %%rax")
	receiverType := receiver.getGtype()
	assert(receiverType.getKind() == G_INTERFACE, nil, "should be interface")

	// dereference: convert an interface value to a concrete value
	receiver.emit()

	emit("LOAD_8_BY_DEREF")

	emit("PUSH_PRIMITIVE # receiver")

	otherArgs := args[1:]
	for i, arg := range otherArgs {
		if _, ok := arg.(*ExprVaArg); ok {
			// skip VaArg for now
			emit("mov $0, %%rax")
		} else {
			arg.emit()
		}
		emit("PUSH_PRIMITIVE # argument no %d", i+2)
	}

	for i, _ := range args {
		j := len(args) - 1 - i
		emit("POP_TO_ARG_%d", j)
	}

	emit("pop %%rax")
	emit("call *%%rax")
}

func loadMapIndexExpr(_map Expr, index Expr) {
	// e.g. x[key]
	emit("# emit map index expr")
	emit("# r10: map header address")
	emit("# r11: map len")
	emit("# r12: specified index value")
	emit("# r13: loop counter")

	// rax: found value (zero if not found)
	// rcx: ok (found: address of the index,  not found:0)
	emit("# emit mapData head address")
	_map.emit()
	emit("mov %%rax, %%r10 # copy head address")
	emitOffsetLoad(_map, IntSize, IntSize)
	emit("mov %%rax, %%r11 # copy len ")
	index.emit()
	emit("mov %%rax, %%r12 # index value")
	emitMapGet(_map.getGtype(), true)
}

func mapOkRegister(is24Width bool) string {
	if is24Width {
		return "rdx"
	} else {
		return "rbx"
	}
}

func emitMapGet(mapType *Gtype, deref bool) {
	if mapType.kind == G_NAMED {
		// @TODO handle infinite chain of relations
		mapType = mapType.relation.gtype
	}
	mapKeyType := mapType.mapKey
	mapValueType := mapType.mapValue
	is24Width := mapValueType.is24Width()
	emit("# emitMapGet")
	emit("mov $0, %%r13 # init loop counter") // i = 0

	labelBegin := makeLabel()
	labelEnd := makeLabel()
	emit("%s: # begin loop ", labelBegin)

	labelIncr := makeLabel()

	emit("cmp %%r11, %%r13") // right, left
	emit("setl %%al # eval(r13(i) < r11(len))")
	emit("movzb %%al, %%eax")
	emit("test %%rax, %%rax")
	if is24Width {
		emit("LOAD_EMPTY_SLICE # NOT FOUND")
	} else if mapValueType.isString() {
		emitEmptyString()
	} else {
		emit("mov $0, %%rax # key not found")
	}

	okRegister := mapOkRegister(is24Width)
	emit("mov $0, %%%s # ok = false", okRegister)

	emit("je %s  # NOT FOUND. exit loop if test makes zero", labelEnd)

	emit("# check if key matches")
	emit("mov %%r13, %%rax")   // i
	emit("imul $16, %%rax")    // i * 16
	emit("mov %%r10, %%rcx")   // head
	emit("add %%rax, %%rcx")   // head + i * 16
	emit("mov (%%rcx), %%rax") // emit index address

	assert(mapKeyType != nil, nil, "key kind should not be nil:"+mapType.String())
	if !mapKeyType.isString() {
		emit("mov (%%rax), %%rax") // dereference
	}
	if mapKeyType.isString() {
		emit("push %%r13")
		emit("push %%r11")
		emit("push %%r10")
		emit("push %%rcx")

		emit("push %%rax")
		emit("push %%r12")
		emitStringsEqualFromStack(true)

		emit("pop %%rcx")
		emit("pop %%r10")
		emit("pop %%r11")
		emit("pop %%r13")
	} else {
		// primitive comparison
		emit("cmp %%r12, %%rax # compare specifiedvalue vs indexvalue")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	}

	emit("test %%rax, %%rax")
	emit("je %s  # Not match. go to next iteration", labelIncr)

	emit("# Value found!")
	emit("mov 8(%%rcx), %%rax # set the found value address")
	emit("mov %%rcx, %%r12 # stash key address")
	if deref {
		if mapValueType.is24Width() {
			emit("mov %%rax, %%r13 # stash")
			emit("mov (%%r13), %%rax # deref 1st")
			emit("mov 8(%%r13), %%rbx # deref 2nd")
			emit("mov 16(%%r13), %%rcx # deref 3rd")
		} else {
			emit("mov (%%rax), %%rax # dereference")
		}
	}

	emit("mov $1, %%%s # ok = true", okRegister)

	emit("jmp %s # exit loop", labelEnd)

	emit("%s: # incr", labelIncr)
	emit("add $1, %%r13") // i++
	emit("jmp %s", labelBegin)

	emit("%s: # end loop", labelEnd)
}

// m[k] = v
// append key and value to the tail of map data, and increment its length
func (e *ExprIndex) emitMapSet(isWidth24 bool) {

	labelAppend := makeLabel()
	labelSave := makeLabel()

	// map get to check if exists
	e.emit()
	// jusdge update or append
	emit("cmp $1, %%%s # ok == true", mapOkRegister(isWidth24))
	emit("sete %%al")
	emit("movzb %%al, %%eax")
	emit("test %%rax, %%rax")
	emit("je %s  # jump to append if not found", labelAppend)

	// update
	emit("push %%r12") // push address of the key
	emit("jmp %s", labelSave)

	// append
	emit("%s: # append to a map ", labelAppend)
	e.collection.emit() // emit pointer address to %rax
	emit("push %%rax # stash head address of mapData")

	// emit len of the map
	elen := &ExprLen{
		arg: e.collection,
	}
	elen.emit()
	emit("imul $%d, %%rax", 2*8) // distance from head to tail
	emit("pop %%rcx")            // head
	emit("add %%rax, %%rcx")     // now rcx is the tail address
	emit("push %%rcx")

	// map len++
	elen.emit()
	emit("add $1, %%rax")
	emitOffsetSave(e.collection, IntSize, ptrSize) // update map len

	// Save key and value
	emit("%s: # end loop", labelSave)
	e.index.emit()
	emit("push %%rax") // index value

	mapType := e.collection.getGtype().Underlying()
	mapKeyType := mapType.mapKey

	if mapKeyType.isString() {
		emit("pop %%rcx")          // index value
		emit("pop %%rax")          // map tail address
		emit("mov %%rcx, (%%rax)") // save indexvalue to malloced area
		emit("push %%rax")         // push map tail
	} else {
		// malloc(8)
		emitCallMalloc(8)
		// %%rax : malloced address
		// stack : [map tail address, index value]
		emit("pop %%rcx")            // index value
		emit("mov %%rcx, (%%rax)")   // save indexvalue to malloced area
		emit("pop %%rcx")            // map tail address
		emit("mov %%rax, (%%rcx) #") // save index address to the tail
		emit("push %%rcx")           // push map tail
	}

	// save value

	// malloc(8)
	var size int = 8
	if isWidth24 {
		size = 24
	}
	emitCallMalloc(size)

	emit("pop %%rcx")           // map tail address
	emit("mov %%rax, 8(%%rcx)") // set malloced address to tail+8

	if isWidth24 {
		emit("pop %%rdx") // rhs value 3/3
		emit("pop %%rcx") // rhs value 2/3
		emit("pop %%rbx") // rhs value 1/3
		// save value
		emit("mov %%rbx, (%%rax)")
		emit("mov %%rcx, 8(%%rax)")
		emit("mov %%rdx, 16(%%rax)")
	} else {
		emit("pop %%rcx") // rhs value
		// save value
		emit("mov %%rcx, (%%rax)") // save value address to the malloced area

	}
}

func (f *StmtFor) emitRangeForMap() {
	emit("# for range %s", f.rng.rangeexpr.getGtype().String())
	assertNotNil(f.rng.indexvar != nil, f.rng.tok)
	labelBegin := makeLabel()
	f.labelEndBlock = makeLabel()
	f.labelEndLoop = makeLabel()

	mapCounter := &Relation{
		name: "",
		expr: f.rng.invisibleMapCounter,
	}
	// counter = 0
	initstmt := &StmtAssignment{
		lefts: []Expr{
			mapCounter,
		},
		rights: []Expr{
			&ExprNumberLiteral{
				val: 0,
			},
		},
	}
	emit("# init index")
	initstmt.emit()

	emit("%s: # begin loop ", labelBegin)

	// counter < len(list)
	condition := &ExprBinop{
		op:   "<",
		left: mapCounter, // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: f.rng.rangeexpr, // len(expr)
		},
	}
	condition.emit()
	emit("test %%rax, %%rax")
	emit("je %s  # if false, exit loop", f.labelEndLoop)

	// set key and value
	mapCounter.emit()
	emit("imul $16, %%rax")
	emit("push %%rax")
	f.rng.rangeexpr.emit() // emit address of map data head
	mapType := f.rng.rangeexpr.getGtype().Underlying()
	mapKeyType := mapType.mapKey

	emit("pop %%rcx")
	emit("add %%rax, %%rcx")
	emit("mov (%%rcx), %%rax")
	if !mapKeyType.isString() {
		emit("mov (%%rax), %%rax")
	}
	f.rng.indexvar.emitSave()

	if f.rng.valuevar != nil {
		emit("# Setting valuevar")
		emit("## rangeexpr.emit()")
		f.rng.rangeexpr.emit()
		emit("mov %%rax, %%rcx # ptr")

		emit("## mapCounter.emit()")
		mapCounter.emit()

		//assert(f.rng.valuevar.getGtype().getSize() <= 8, f.rng.token(), "invalid size")
		emit("## eval value")
		emit("imul $16, %%rax  # counter * 16")
		emit("add $8, %%rax    # counter * 16 + 8")
		emit("add %%rax, %%rcx # mapHead + (counter * 16 + 8)")
		emit("mov (%%rcx), %%rdx")

		switch f.rng.valuevar.getGtype().getKind() {
		case G_SLICE, G_MAP:
			emit("mov (%%rdx), %%rax")
			emit("mov 8(%%rdx), %%rbx")
			emit("mov 16(%%rdx), %%rcx")
			emit("push %%rax")
			emit("push %%rbx")
			emit("push %%rcx")
			emitSave24(f.rng.valuevar, 0)
		default:
			emit("mov (%%rdx), %%rax")
			f.rng.valuevar.emitSave()
		}

	}

	f.block.emit()
	emit("%s: # end block", f.labelEndBlock)

	// counter++
	indexIncr := &StmtInc{
		operand: mapCounter,
	}
	indexIncr.emit()

	emit("jmp %s", labelBegin)
	emit("%s: # end loop", f.labelEndLoop)
}

// push addr, len, cap
func (lit *ExprMapLiteral) emit() {
	length := len(lit.elements)

	// allocaated address of the map head
	var size int
	if length == 0 {
		size = ptrSize * 1024
	} else {
		size = length * ptrSize * 1024
	}
	emitCallMalloc(size)
	emit("push %%rax") // map head

	mapType := lit.getGtype()
	mapKeyType := mapType.mapKey

	for i, element := range lit.elements {
		// alloc key
		if mapKeyType.isString() {
			element.key.emit()
		} else {
			element.key.emit()
			emit("push %%rax") // value of key
			// call malloc for key
			emitCallMalloc(8)
			emit("pop %%rcx")          // value of key
			emit("mov %%rcx, (%%rax)") // save key to heap
		}

		emit("pop %%rbx")                     // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8) // save key address
		emit("push %%rbx")                    // map head

		if element.value.getGtype().getSize() <= 8 {
			element.value.emit()
			emit("push %%rax") // value of value
			// call malloc
			emitCallMalloc(8)
			emit("pop %%rcx")          // value of value
			emit("mov %%rcx, (%%rax)") // save value to heap
		} else {
			switch element.value.getGtype().getKind() {
			case G_MAP, G_SLICE, G_INTERFACE:
				// rax,rbx,rcx
				element.value.emit()
				emit("push %%rax") // ptr
				emitCallMalloc(8 * 3)
				emit("pop %%rdx") // ptr
				emit("mov %%rdx, %d(%%rax)", 0)
				emit("mov %%rbx, %d(%%rax)", 8)
				emit("mov %%rcx, %d(%%rax)", 16)

			default:
				TBI(element.value.token(), "unable to handle %s", element.value.getGtype())
			}
		}

		emit("pop %%rbx") // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8+8)
		emit("push %%rbx")
	}

	emit("pop %%rax")
	emit("push %%rax")       // address (head of the heap)
	emit("push $%d", length) // len
	emit("push $%d", length) // cap

	emit("POP_MAP")
}
