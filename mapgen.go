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
	emit("TEST_IT")
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
	emit("IMUL_NUMBER 16")    // i * 16
	emit("PUSH_PRIMITIVE")

	emit("mov %%r10, %%rax")   // head
	emit("PUSH_PRIMITIVE")

	emit("SUM_FROM_STACK")   // head + i * 16

	emit("PUSH_PRIMITIVE") // index address
	emit("LOAD_8_BY_DEREF") // emit index address

	assert(mapKeyType != nil, nil, "key kind should not be nil:"+mapType.String())
	if !mapKeyType.isString() {
		emit("LOAD_8_BY_DEREF") // dereference
	}
	if mapKeyType.isString() {
		emit("push %%r13")
		emit("push %%r11")
		emit("push %%r10")

		emit("push %%rax")
		emit("push %%r12")
		emitStringsEqualFromStack(true)

		emit("pop %%r10")
		emit("pop %%r11")
		emit("pop %%r13")
	} else {
		// primitive comparison
		emit("cmp %%r12, %%rax # compare specifiedvalue vs indexvalue")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	}
	emit("pop %%rcx") // index address

	emit("TEST_IT")
	emit("je %s  # Not match. go to next iteration", labelIncr)

	emit("# Value found!")
	emit("mov 8(%%rcx), %%rax # set the found value address")
	emit("mov %%rcx, %%r12 # stash key address")
	if deref {
		if mapValueType.is24Width() {
			emit("mov %%rax, %%r13 # stash")
			emit("LOAD_24_BY_DEREF")
		} else {
			emit("LOAD_8_BY_DEREF")
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
	emit("TEST_IT")
	emit("je %s  # jump to append if not found", labelAppend)

	// update
	emit("push %%r12") // push address of the key
	emit("jmp %s", labelSave)

	// append
	emit("%s: # append to a map ", labelAppend)
	e.collection.emit() // emit pointer address to %rax
	emit("PUSH_PRIMITIVE")

	// emit len of the map
	elen := &ExprLen{
		arg: e.collection,
	}
	elen.emit()
	emit("IMUL_NUMBER %d", 2*8) // distance from head to tail
	emit("PUSH_PRIMITIVE")
	emit("SUM_FROM_STACK")
	emit("PUSH_PRIMITIVE")

	// map len++
	elen.emit()
	emit("ADD_NUMBER 1")
	emitOffsetSave(e.collection, IntSize, ptrSize) // update map len

	// Save key and value
	emit("%s: # end loop", labelSave)
	e.index.emit()
	emit("PUSH_PRIMITIVE") // index value

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
	emit("PUSH_PRIMITIVE")
	if isWidth24 {
		emit("STORE_24_INDIRECT_FROM_STACK")
	} else {
		emit("STORE_8_INDIRECT_FROM_STACK")
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
	emit("TEST_IT")
	emit("je %s  # if false, exit loop", f.labelEndLoop)

	// set key and value
	mapCounter.emit()
	emit("IMUL_NUMBER 16")
	emit("PUSH_PRIMITIVE # x")
	f.rng.rangeexpr.emit() // emit address of map data head
	emit("PUSH_PRIMITIVE # y")

	mapType := f.rng.rangeexpr.getGtype().Underlying()
	mapKeyType := mapType.mapKey

	emit("SUM_FROM_STACK # x + y")
	emit("LOAD_8_BY_DEREF")

	if !mapKeyType.isString() {
		emit("LOAD_8_BY_DEREF")
	}
	f.rng.indexvar.emitSave()

	if f.rng.valuevar != nil {
		emit("# Setting valuevar")
		emit("## rangeexpr.emit()")
		f.rng.rangeexpr.emit()
		emit("PUSH_PRIMITIVE")

		emit("## mapCounter.emit()")
		mapCounter.emit()
		emit("## eval value")
		emit("IMUL_NUMBER 16  # counter * 16")
		emit("ADD_NUMBER 8 # counter * 16 + 8")
		emit("PUSH_PRIMITIVE")

		emit("SUM_FROM_STACK")

		emit("LOAD_8_BY_DEREF")

		switch f.rng.valuevar.getGtype().getKind() {
		case G_SLICE, G_MAP:
			emit("LOAD_24_BY_DEREF")
			emit("PUSH_24")
			emitSave24(f.rng.valuevar, 0)
		default:
			emit("LOAD_8_BY_DEREF")
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
			emit("PUSH_PRIMITIVE") // value of key
			// call malloc for key
			emitCallMalloc(8)
			emit("PUSH_PRIMITIVE")

			emit("STORE_8_INDIRECT_FROM_STACK") // save key to heap
		}

		emit("pop %%rbx")                     // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8) // save key address
		emit("push %%rbx")                    // map head

		if element.value.getGtype().getSize() <= 8 {
			element.value.emit()
			emit("push %%rax") // value of value
			emitCallMalloc(8)
			emit("PUSH_PRIMITIVE")
			emit("STORE_8_INDIRECT_FROM_STACK") // save value to heap
		} else {
			switch element.value.getGtype().getKind() {
			case G_MAP, G_SLICE, G_INTERFACE:
				// rax,rbx,rcx
				element.value.emit()
				emit("PUSH_24") // ptr
				emitCallMalloc(8 * 3)
				emit("PUSH_PRIMITIVE")
				emit("STORE_24_INDIRECT_FROM_STACK")
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
