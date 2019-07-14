package main

const MAX_METHODS_PER_TYPE int = 128

func (call *IrInterfaceMethodCall) emit() {
	receiver := call.receiver
	var methodName bytes = bytes(call.methodName)
	emit("# emit interface method call \"%s\"", methodName)
	mapType := &Gtype{
		kind: G_MAP,
		mapKey: &Gtype{
			kind: G_POINTER,
		},
		mapValue: &Gtype{
			kind: G_POINTER,
		},
	}
	emit("# emit receiverTypeId of %s", receiver.getGtype().String())
	emitOffsetLoad(receiver, ptrSize, ptrSize)
	emit("IMUL_NUMBER 8")
	emit("PUSH_8")

	emit("lea receiverTypes(%%rip), %%rax")
	emit("PUSH_8")
	emit("SUM_FROM_STACK")

	emit("# find method %s", methodName)
	emit("mov (%%rax), %%rax") // address of receiverType
	emit("PUSH_8 # map head")

	emit("LOAD_NUMBER %d", MAX_METHODS_PER_TYPE) // max methods for a type
	emit("PUSH_8 # len")

	emit("lea .S.%s, %%rax", methodName) // index value (addr)
	emit("PUSH_8 # map index value")

	emitMapGet(mapType, false)

	emit("PUSH_8 # funcref")

	emit("mov $0, %%rax")
	receiverType := receiver.getGtype()
	assert(receiverType.getKind() == G_INTERFACE, nil, S("should be interface"))

	receiver.emit()
	emit("LOAD_8_BY_DEREF # dereference: convert an interface value to a concrete value")

	emit("PUSH_8 # receiver")

	call.emitMethodCall()
}

func emitLoadMapKey(eMapKey Expr) {
	var isKeyString bool = eMapKey.getGtype().isString()
	if isKeyString {
		var arg0 Expr
		switch eMapKey.(type) {
		case *ExprFuncallOrConversion:
			funcallOrConversion := eMapKey.(*ExprFuncallOrConversion)
			arg0 = funcallOrConversion.args[0]
		case *IrExprConversion:
			conversion := eMapKey.(*IrExprConversion)
			arg0 = conversion.arg
		default:
			assertNotReached(eMapKey.token())
		}
		arg0.emit()
	} else {
		eMapKey.emit()
	}
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
	emit("# not nil")
	emit("LOAD_8_BY_DEREF")
	emit("PUSH_8 # map head")
	_map.emit()
	emit("mov 8(%%rax), %%rax")
	emit("PUSH_8 # len")
	var isKeyString bool = e.index.getGtype().isString()
	emitLoadMapKey(e.index)
	if isKeyString {
		emit("PUSH_SLICE")
		emit("pop %%rdi # index value")
		emit("pop %%rdx # index value")
		emit("pop %%rcx # index value")
	} else {
		emit("PUSH_8 # index value")
		emit("pop %%rcx # index value")
	}

	emit("pop %%rbx # len")
	emit("pop %%rax # heap")

	emit("jmp %s", labelEnd)
	// nil case
	emit("%s:", labelNil)
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emit("mov $0, %%rdx")
	emit("mov $0, %%rdi")
	emit("%s:", labelEnd)

	emit("push %%rax")
	emit("push %%rbx")
	if isKeyString {
		emit("push %%rcx")
		emit("push %%rdx")
		emit("push %%rdi")
	} else {
		emit("push %%rcx")
	}
	emitMapGet(_map.getGtype(), isKeyString)
}

func mapOkRegister(is24Width bool) bytes {
	if is24Width {
		return S("rdx")
	} else {
		return S("rbx")
	}
}

// get map from stack
// r10: map header address")
// r11: map len")
// r12: loop counter")
// r13: k value in m[k]
// r14: k value in m[k]
// r15: k value in m[k]
func emitMapGet(mapType *Gtype, cmpGoString bool) {

	mapType = mapType.Underlying()
	mapKeyType := mapType.mapKey
	mapValueType := mapType.mapValue
	isValue24Width := mapValueType.is24WidthType()
	labelBegin := makeLabel()
	labelEnd := makeLabel()
	labelIncr := makeLabel()

	emit("# emitMapGet")

	if cmpGoString {
		emit("pop %%r15") // index value a
		emit("pop %%r14") // index value b
		emit("pop %%r13") // index value c
	} else {
		emit("pop %%r13") // index value
	}

	emit("pop %%r11") // len
	emit("pop %%r10") // head addr

	emit("mov $0, %%r12 # init loop counter") // i = 0

	emit("%s: # begin loop ", labelBegin)

	emit("push %%r12 # loop counter")
	emit("push %%r11 # map len")
	emit("CMP_FROM_STACK setl")
	emit("TEST_IT")
	if isValue24Width {
		emit("LOAD_EMPTY_SLICE # NOT FOUND")
	} else {
		emit("mov $0, %%rax # key not found")
	}

	okRegister := mapOkRegister(isValue24Width)
	emit("mov $0, %%%s # ok = false", okRegister)

	emit("je %s  # Exit. NOT FOUND IN ALL KEYS.", labelEnd)

	emit("# check if key matches")
	emit("mov %%r12, %%rax") // i
	emit("IMUL_NUMBER 16")   // i * 16
	emit("PUSH_8")

	emit("mov %%r10, %%rax") // head
	emit("PUSH_8")

	emit("SUM_FROM_STACK") // head + i * 16

	emit("PUSH_8")          // addr of key addr
	emit("LOAD_8_BY_DEREF") // emit key addr

	assert(mapKeyType != nil, nil, S("key kind should not be nil:%s"), mapType.String())

	if cmpGoString {
		emit("push %%r12")
		emit("push %%r11")
		emit("push %%r10")

		emit("LOAD_24_BY_DEREF") // dereference
		emit("PUSH_SLICE")

		// push index value
		emit("push %%r13")
		emit("push %%r14")
		emit("push %%r15")

		emitGoStringsEqualFromStack()

		emit("pop %%r10")
		emit("pop %%r11")
		emit("pop %%r12")
	} else {
		emit("LOAD_8_BY_DEREF") // dereference
		// primitive comparison
		emit("cmp %%r13, %%rax # compare specifiedvalue vs indexvalue")
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
	if mapValueType.is24WidthType() {
		emit("LOAD_24_BY_DEREF")
	} else {
		emit("LOAD_8_BY_DEREF")
	}

	emit("mov $1, %%%s # ok = true", okRegister)
	emit("pop %%r13 # key address. will be in map set")
	emit("jmp %s # exit loop", labelEnd)

	emit("%s: # incr", labelIncr)
	emit("add $1, %%r12") // i++
	emit("jmp %s", labelBegin)

	emit("%s: # end loop", labelEnd)

}

// m[k] = v
func (e *ExprIndex) emitMapSetFromStack8() {
	e.emitMapSetFromStack(false)
}

// m[k] = v
func (e *ExprIndex) emitMapSetFromStack24() {
	e.emitMapSetFromStack(true)
}

func emitSaveMapKey(eMapKey Expr) {
	emitLoadMapKey(eMapKey)
	if eMapKey.getGtype().isString() {
		emit("PUSH_SLICE")
		emitCallMalloc(24)
		emit("PUSH_8")
		emit("STORE_24_INDIRECT_FROM_STACK") // save key to mallocedaddr
	} else {
		emit("PUSH_8") // index value
		emitCallMalloc(8)
		emit("PUSH_8")
		emit("STORE_8_INDIRECT_FROM_STACK") // save key to mallocedaddr
	}

	// %rax is mappedaddr
}

// m[k] = v
// append key and value to the tail of map data, and increment its length
func (e *ExprIndex) emitMapSetFromStack(isValueWidth24 bool) {
	emit("# emitMapSetFromStack")
	labelAppend := makeLabel()
	labelSave := makeLabel()

	// map get to check if exists
	e.emit()
	// jusdge update or append
	emit("cmp $1, %%%s # ok == true", mapOkRegister(isValueWidth24))
	emit("sete %%al")
	emit("movzb %%al, %%eax")
	emit("TEST_IT")
	emit("je %s  # jump to append if not found", labelAppend)

	// update
	emit("push %%r13") // push address of the key
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
	emit("IMUL_NUMBER %d", mapUnitSize) // distance from head to tail
	emit("PUSH_8")
	emit("SUM_FROM_STACK")
	emit("PUSH_8") // tail addr

	// map len++
	elen.emit()
	emit("ADD_NUMBER 1")
	emit("PUSH_8")
	e.collection.emit()
	emit("pop %%rbx # new len")
	emit("mov %%rbx, 8(%%rax) # update map len")

	// Save key and value
	emit("%s: # end loop", labelSave)

	// save key
	emit("# save key")

	emitSaveMapKey(e.index)
	emit("mov %%rax, %%rcx")   // copy mallocedaddr
	// append key to tail
	emit("POP_8")              // tailaddr
	emit("mov %%rcx, (%%rax)") // save mallocedaddr to tailaddr
	emit("PUSH_8")             // push tailaddr


	// save value

	// malloc(8)
	var size int = 8
	if isValueWidth24 {
		size = 24
	}
	emitCallMalloc(size)

	emit("pop %%rcx")           // map tail address
	emit("mov %%rax, 8(%%rcx)") // set malloced address to tail+8
	emit("PUSH_8")
	if isValueWidth24 {
		emit("STORE_24_INDIRECT_FROM_STACK")
	} else {
		emit("STORE_8_INDIRECT_FROM_STACK")
	}
}

func (em *IrStmtRangeMap) emit() {
	//mapType := em.rangeexpr.getGtype().Underlying()
	//mapKeyType := mapType.mapKey

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
		op:   bytes("<"),
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

	if em.indexvar.getGtype().isString() {
		emit("LOAD_24_BY_DEREF")
		emitSave24(em.indexvar,0)
	} else {
		emit("LOAD_8_BY_DEREF")
		emitSavePrimitive(em.indexvar)
	}

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

var mapUnitSize int = 2*8

// push addr, len, cap
func (lit *ExprMapLiteral) emit() {
	var length int = len(lit.elements)

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

	//mapType := lit.getGtype()
	//mapKeyType := mapType.mapKey

	for i, element := range lit.elements {
		var offsetKey int = i*mapUnitSize
		var offsetValue int = i*mapUnitSize+8

		// save key
		emitSaveMapKey(element.key)
		emit("mov %%rax, %%rcx")   // copy mallocedaddr
		// append key to tail
		emit("POP_8")                         // map head
		emit("mov %%rcx, %d(%%rax) #", offsetKey) // save key address
		emit("PUSH_8")                        // map head

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
			TBI(element.value.token(), S("unable to handle %s"), element.value.getGtype())

		}

		emit("pop %%rbx") // map head
		emit("mov %%rax, %d(%%rbx) #", offsetValue)
		emit("push %%rbx")
	}

	emitCallMalloc(16)
	emit("pop %%rbx") // address (head of the heap)
	emit("mov %%rbx, (%%rax)")
	emit("mov $%d, %%rcx", length) // len
	emit("mov %%rcx, 8(%%rax)")
}
