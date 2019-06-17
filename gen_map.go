package main

const MAX_METHODS_PER_TYPE int = 128

func (call *IrInterfaceMethodCall) emit() {
	receiver := call.receiver
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
	emit("# emit receiverTypeId of %s", receiver.getGtype().String())
	emitOffsetLoad(receiver, ptrSize, ptrSize)
	emit("IMUL_NUMBER 8")
	emit("PUSH_8")

	emit("lea receiverTypes(%%rip), %%rax")
	emit("PUSH_8")
	emit("SUM_FROM_STACK")

	emit("# find method %s", call.methodName)
	emit("mov (%%rax), %%rax") // address of receiverType
	emit("PUSH_8 # map head")

	emit("LOAD_NUMBER %d", MAX_METHODS_PER_TYPE) // max methods for a type
	emit("PUSH_8 # len")

	emit("lea .M%s, %%rax", call.methodName) // index value
	emit("PUSH_8 # map index value")                 // index value

	emit("pop %%rcx")
	emit("pop %%rbx")
	emit("pop %%rax")
	emitMapGet(mapType, false)

	emit("PUSH_8")

	emit("mov $0, %%rax")
	receiverType := receiver.getGtype()
	assert(receiverType.getKind() == G_INTERFACE, nil, "should be interface")

	// dereference: convert an interface value to a concrete value
	receiver.emit()

	emit("LOAD_8_BY_DEREF")

	emit("PUSH_8 # funcref")

	call.emitMethodCall()
}

// emit map index expr
func loadMapIndexExpr(e *ExprIndex) {
	// e.g. x[key]

	_map := e.collection
	// rax: found value (zero if not found)
	// rcx: ok (found: address of the index,  not found:0)
	emit("# emit mapData head address")
	_map.emit()

	// if not nil
	// then emit 24width data
	// else emit 24width 0
	labelNil := makeLabel()
	labelEnd := makeLabel()
	emit("TEST_IT # map && map (check if map is nil)")
	emit("je %s # jump if map is nil", labelNil)
	// not nil case
	emit("# not null")
	emit("LOAD_8_BY_DEREF")
	emit("PUSH_8 # map head")
	_map.emit()
	emit("mov 8(%%rax), %%rax")
	emit("PUSH_8 # len")
	e.index.emit()
	emit("PUSH_8 # index value")

	emit("pop %%rcx # index value")
	emit("pop %%rbx # len")
	emit("pop %%rax # heap")

	emit("jmp %s", labelEnd)
	// nil case
	emit("%s:", labelNil)
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emit("%s:", labelEnd)


	emitMapGet(_map.getGtype(), true)
}

func mapOkRegister(is24Width bool) string {
	if is24Width {
		return "rdx"
	} else {
		return "rbx"
	}
}

// r10: map header address")
// r11: map len")
// r12: specified index value")
// r13: loop counter")
func emitMapGet(mapType *Gtype, deref bool) {

	mapType = mapType.Underlying()
	mapKeyType := mapType.mapKey
	mapValueType := mapType.mapValue
	is24Width := mapValueType.is24WidthType()

	emit("# emitMapGet")

	emit("mov %%rax, %%r10")
	emit("mov %%rbx, %%r11")
	emit("mov %%rcx, %%r12")

	labelBegin := makeLabel()
	labelEnd := makeLabel()
	labelIncr := makeLabel()

	emit("mov $0, %%r13 # init loop counter") // i = 0

	emit("%s: # begin loop ", labelBegin)

	emit("push %%r13")
	emit("push %%r11")
	emit("CMP_FROM_STACK setl")
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

	emit("je %s  # Exit. NOT FOUND IN ALL KEYS.", labelEnd)

	emit("# check if key matches")
	emit("mov %%r13, %%rax")   // i
	emit("IMUL_NUMBER 16")    // i * 16
	emit("PUSH_8")

	emit("mov %%r10, %%rax")   // head
	emit("PUSH_8")

	emit("SUM_FROM_STACK")   // head + i * 16

	emit("PUSH_8") // index address
	emit("LOAD_8_BY_DEREF") // emit index address

	assert(mapKeyType != nil, nil, "key kind should not be nil:"+mapType.String())
	if !mapKeyType.isString() {
		emit("LOAD_8_BY_DEREF") // dereference
	}
	if mapKeyType.isString() {
		emit("push %%r13")
		emit("push %%r11")
		emit("push %%r10")

		emit("PUSH_8")
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

	emit("TEST_IT")
	emit("pop %%rax") // index address
	emit("je %s  # Not match. go to next iteration", labelIncr)

	emit("# Value found!")
	emit("push %%rax # stash key address")
	emit("ADD_NUMBER 8 # value address")
	emit("LOAD_8_BY_DEREF # set the found value address")
	if deref {
		if mapValueType.is24WidthType() {
			emit("LOAD_24_BY_DEREF")
		} else {
			emit("LOAD_8_BY_DEREF")
		}
	}

	emit("mov $1, %%%s # ok = true", okRegister)
	emit("pop %%r12 # key address. will be in map set")
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
	emit("LOAD_8_BY_DEREF")
	emit("PUSH_8")

	// emit len of the map
	elen := &ExprLen{
		arg: e.collection,
	}
	elen.emit()
	emit("IMUL_NUMBER %d", 2*8) // distance from head to tail
	emit("PUSH_8")
	emit("SUM_FROM_STACK")
	emit("PUSH_8")

	// map len++
	elen.emit()
	emit("ADD_NUMBER 1")
	emit("PUSH_8")
	e.collection.emit()
	emit("pop %%rbx # new len")
	emit("mov %%rbx, 8(%%rax) # update map len")

	// Save key and value
	emit("%s: # end loop", labelSave)
	e.index.emit()
	emit("PUSH_8") // index value

	mapType := e.collection.getGtype().Underlying()
	mapKeyType := mapType.mapKey

	if mapKeyType.isString() {
		emit("pop %%rcx")          // index value
		emit("pop %%rax")          // map tail address
		emit("mov %%rcx, (%%rax)") // save indexvalue to malloced area
		emit("PUSH_8")         // push map tail
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
	emit("PUSH_8")
	if isWidth24 {
		emit("STORE_24_INDIRECT_FROM_STACK")
	} else {
		emit("STORE_8_INDIRECT_FROM_STACK")
	}
}

type RangeMapEmitter struct {
	tok *Token
	labels *LoopLabels
	rangeexpr Expr
	indexvar *Relation
	valuevar *Relation
	mapCounter *ExprVariable
	initstmt Stmt
	condition Expr
	indexIncr Stmt
	block Stmt
}

func (em *RangeMapEmitter) token() *Token {
	return em.tok
}

func (em *RangeMapEmitter) emit() {
	mapType := em.rangeexpr.getGtype().Underlying()
	mapKeyType := mapType.mapKey

	// counter = 0
	em.initstmt = &StmtAssignment{
		lefts: []Expr{
			em.mapCounter,
		},
		rights: []Expr{
			&ExprNumberLiteral{
				val: 0,
			},
		},
	}
	// counter < len(list)
	em.condition = &ExprBinop{
		op:   "<",
		left: em.mapCounter, // i
		// @TODO
		// The range expression x is evaluated once before beginning the loop
		right: &ExprLen{
			arg: em.rangeexpr, // len(expr)
		},
	}

	// counter++
	em.indexIncr = &StmtInc{
		operand: em.mapCounter,
	}

	emit("# init index")
	em.initstmt.emit()

	emit("%s: # begin loop ", em.labels.labelBegin)

	em.condition.emit()
	emit("TEST_IT")
	emit("je %s  # if false, exit loop", em.labels.labelEndLoop)

	// set key and value
	em.mapCounter.emit()
	emit("IMUL_NUMBER 16")
	emit("PUSH_8 # x")
	em.rangeexpr.emit() // emit address of map data head
	emit("LOAD_8_BY_DEREF")
	emit("PUSH_8 # y")

	emit("SUM_FROM_STACK # x + y")
	emit("LOAD_8_BY_DEREF")

	if !mapKeyType.isString() {
		emit("LOAD_8_BY_DEREF")
	}
	emitSavePrimitive(em.indexvar)

	if em.valuevar != nil {
		emit("# Setting valuevar")
		emit("## rangeexpr.emit()")
		em.rangeexpr.emit()
		emit("LOAD_8_BY_DEREF# map head")
		emit("PUSH_8")

		emit("## mapCounter.emit()")
		em.mapCounter.emit()
		emit("## eval value")
		emit("IMUL_NUMBER 16  # counter * 16")
		emit("ADD_NUMBER 8 # counter * 16 + 8")
		emit("PUSH_8")

		emit("SUM_FROM_STACK")

		emit("LOAD_8_BY_DEREF")

		if em.valuevar.getGtype().is24WidthType() {
			emit("LOAD_24_BY_DEREF")
			emitSave24(em.valuevar, 0)
		} else {
			emit("LOAD_8_BY_DEREF")
			emitSavePrimitive(em.valuevar)
		}

	}

	em.block.emit()
	emit("%s: # end block", em.labels.labelEndBlock)

	em.indexIncr.emit()

	emit("jmp %s", em.labels.labelBegin)
	emit("%s: # end loop", em.labels.labelEndLoop)
}

// push addr, len, cap
func (lit *ExprMapLiteral) emit() {
	length := len(lit.elements)

	// allocaated address of the map head
	// @FIXME 1024 is a tentative number
	var size int
	if length == 0 {
		size = ptrSize * 1024
	} else {
		size = length * ptrSize * 1024
	}
	emitCallMalloc(size)
	emit("PUSH_8") // map head

	mapType := lit.getGtype()
	mapKeyType := mapType.mapKey

	for i, element := range lit.elements {
		// alloc key
		if mapKeyType.isString() {
			element.key.emit()
		} else {
			element.key.emit()
			emit("PUSH_8") // value of key
			// call malloc for key
			emitCallMalloc(8)
			emit("PUSH_8")

			emit("STORE_8_INDIRECT_FROM_STACK") // save key to heap
		}

		emit("pop %%rbx")                     // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8) // save key address
		emit("push %%rbx")                    // map head

		if element.value.getGtype().getSize() <= 8 {
			element.value.emit()
			emit("PUSH_8") // value of value
			emitCallMalloc(8)
			emit("PUSH_8")
			emit("STORE_8_INDIRECT_FROM_STACK") // save value to heap
		} else if element.value.getGtype().is24WidthType() {
			// rax,rbx,rcx
			element.value.emit()
			emit("PUSH_24") // ptr
			emitCallMalloc(8 * 3)
			emit("PUSH_8")
			emit("STORE_24_INDIRECT_FROM_STACK")
		} else {
			TBI(element.value.token(), "unable to handle %s", element.value.getGtype())

		}

		emit("pop %%rbx") // map head
		emit("mov %%rax, %d(%%rbx) #", i*2*8+8)
		emit("push %%rbx")
	}

	emitCallMalloc(16)
	emit("pop %%rbx") // address (head of the heap)
	emit("mov %%rbx, (%%rax)")
	emit("mov $%d, %%rcx", length) // len
	emit("mov %%rcx, 8(%%rax)")
}
