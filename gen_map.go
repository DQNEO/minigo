package main

const MAX_METHODS_PER_TYPE int = 128

type IrMapInitializer struct {
	tok      *Token
	gtype    *Gtype
	elements []*MapElement // for map literal
	lenArg   Expr          // for make(T, n)
}

func (e *IrMapInitializer) token() *Token {
	return e.tok
}

func (e *IrMapInitializer) dump() {
	panic("implement me")
}

func (e *IrMapInitializer) getGtype() *Gtype {
	return e.gtype
}

func (call *IrInterfaceMethodCall) emit() {
	receiver := call.receiver
	methodName := call.methodName
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

	emit("leaq receiverTypes(%%rip), %%rax")
	emit("PUSH_8")
	emit("SUM_FROM_STACK")

	emit("# find method %s", methodName)
	emit("movq (%%rax), %%rax") // address of receiverType
	emit("PUSH_8 # map head")

	emit("LOAD_NUMBER %d", MAX_METHODS_PER_TYPE) // max methods for a type
	emit("PUSH_8 # len")

	emit("leaq .S.%s, %%rax", methodName) // index value (addr)
	emit("PUSH_8 # map index value")

	emitMapGet(mapType, false)

	emit("PUSH_8 # funcref")

	emit("movq $0, %%rax")
	receiverType := receiver.getGtype()
	assert(receiverType.getKind() == G_INTERFACE, nil, "should be interface")

	_call := &IrCall{
		isInterfaceMethodCall: true,
		symbol:                "",
		icallee:                call.callee,
		receiver:              call.receiver,
		args:                  call.args,
	}
	_call.emit()
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
	emit("movq 8(%%rax), %%rax")
	emit("PUSH_8 # len")
	var isKeyString bool = e.index.getGtype().isString()
	e.index.emit()
	if isKeyString {
		emit("PUSH_SLICE")
		emit("popq %%rdi # index value")
		emit("popq %%rdx # index value")
		emit("popq %%rcx # index value")
	} else {
		emit("PUSH_8 # index value")
		emit("popq %%rcx # index value")
	}

	emit("popq %%rbx # len")
	emit("popq %%rax # heap")

	emit("jmp %s", labelEnd)
	// nil case
	emit("%s:", labelNil)
	emit("movq $0, %%rax")
	emit("movq $0, %%rbx")
	emit("movq $0, %%rcx")
	emit("movq $0, %%rdx")
	emit("movq $0, %%rdi")
	emit("%s:", labelEnd)

	emit("pushq %%rax")
	emit("pushq %%rbx")
	if isKeyString {
		emit("pushq %%rcx")
		emit("pushq %%rdx")
		emit("pushq %%rdi")
	} else {
		emit("pushq %%rcx")
	}
	emitMapGet(_map.getGtype(), isKeyString)
}

func mapOkRegister(is24Width bool) string {
	if is24Width {
		return "rdx"
	} else {
		return "rbx"
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
		emit("popq %%r15") // index value a
		emit("popq %%r14") // index value b
		emit("popq %%r13") // index value c
	} else {
		emit("popq %%r13") // index value
	}

	emit("popq %%r11") // len
	emit("popq %%r10") // head addr

	emit("movq $0, %%r12 # init loop counter") // i = 0

	emit("%s: # begin loop ", labelBegin)

	emit("pushq %%r12 # loop counter")
	emit("pushq %%r11 # map len")
	emit("CMP_FROM_STACK setl")
	emit("TEST_IT")
	if isValue24Width {
		emit("LOAD_EMPTY_SLICE # NOT FOUND")
	} else {
		emit("movq $0, %%rax # key not found")
	}

	okRegister := mapOkRegister(isValue24Width)
	emit("movq $0, %%%s # ok = false", okRegister)

	emit("je %s  # Exit. NOT FOUND IN ALL KEYS.", labelEnd)

	emit("# check if key matches")
	emit("movq %%r12, %%rax") // i
	emit("IMUL_NUMBER 16")   // i * 16
	emit("PUSH_8")

	emit("movq %%r10, %%rax") // head
	emit("PUSH_8")

	emit("SUM_FROM_STACK") // head + i * 16

	emit("PUSH_8")          // addr of key addr
	emit("LOAD_8_BY_DEREF") // emit key addr

	assert(mapKeyType != nil, nil, "key kind should not be nil:%s", mapType.String())

	if cmpGoString {
		emit("pushq %%r12")
		emit("pushq %%r11")
		emit("pushq %%r10")

		emit("LOAD_24_BY_DEREF") // dereference
		emit("PUSH_SLICE")

		// push index value
		emit("pushq %%r13")
		emit("pushq %%r14")
		emit("pushq %%r15")

		emitGoStringsEqualFromStack()

		emit("popq %%r10")
		emit("popq %%r11")
		emit("popq %%r12")
	} else {
		emit("LOAD_8_BY_DEREF") // dereference
		// primitive comparison
		emit("cmpq %%r13, %%rax # compare specifiedvalue vs indexvalue")
		emit("sete %%al")
		emit("movzb %%al, %%eax")
	}

	emit("TEST_IT")
	emit("popq %%rax") // index address
	emit("je %s  # Not match. go to next iteration", labelIncr)

	emit("# Value found!")
	emit("pushq %%rax # stash key address")
	emit("ADD_NUMBER 8 # value address")
	emit("LOAD_8_BY_DEREF # set the found value address")
	if mapValueType.is24WidthType() {
		emit("LOAD_24_BY_DEREF")
	} else {
		emit("LOAD_8_BY_DEREF")
	}

	emit("movq $1, %%%s # ok = true", okRegister)
	emit("popq %%r13 # key address. will be in map set")
	emit("jmp %s # exit loop", labelEnd)

	emit("%s: # incr", labelIncr)
	emit("addq $1, %%r12") // i++
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
	eMapKey.emit()
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
	emit("cmpq $1, %%%s # ok == true", mapOkRegister(isValueWidth24))
	emit("sete %%al")
	emit("movzb %%al, %%eax")
	emit("TEST_IT")
	emit("je %s  # jump to append if not found", labelAppend)

	// update
	emit("pushq %%r13") // push address of the key
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
	emit("popq %%rbx # new len")
	emit("movq %%rbx, 8(%%rax) # update map len")

	// Save key and value
	emit("%s: # end loop", labelSave)

	// save key
	emit("# save key")

	emitSaveMapKey(e.index)
	emit("movq %%rax, %%rcx") // copy mallocedaddr
	// append key to tail
	emit("POP_8")              // tailaddr
	emit("movq %%rcx, (%%rax)") // save mallocedaddr to tailaddr
	emit("PUSH_8")             // push tailaddr

	// save value

	// malloc(8)
	var size int = 8
	if isValueWidth24 {
		size = 24
	}
	emitCallMalloc(size)

	emit("popq %%rcx")           // map tail address
	emit("movq %%rax, 8(%%rcx)") // set malloced address to tail+8
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

	if em.indexvar.getGtype().isString() {
		emit("LOAD_24_BY_DEREF")
		emitSave24(em.indexvar, 0)
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

var mapUnitSize int = 2 * 8

func (e *ExprMapLiteral) emit() {
	mapInitializer := &IrMapInitializer{
		tok:      e.token(),
		gtype:    e.getGtype(),
		elements: e.elements,
	}
	mapInitializer.emit()
}
// push addr, len, cap
func (e *IrMapInitializer) emit() {
	var multiple int = 1024 // Forgot what it is ...
	if e.lenArg != nil {
		// make(T, n)
		eLen := &ExprBinop{
			tok:   e.token(),
			op:    "*",
			left:  e.lenArg,
			right: &ExprNumberLiteral{
				val: ptrSize * multiple,
			},
		}
		emit("# make map")
		emit("movq $0, %%rax")
		emit("PUSH_8") // map len
		emit("# alloc map area")
		emitCallMallocDinamicSize(eLen)
	} else {
		// make(T) or map[T1]T2{...}
		var length int = len(e.elements)
		// allocaated address of the map head
		var size int
		if length == 0 {
			size = ptrSize * multiple
		} else {
			size = length * (ptrSize * multiple)
		}
		elen := &ExprNumberLiteral{val:length}
		elen.emit()
		emit("PUSH_8") // map len
		emitCallMalloc(size)
	}

	emit("PUSH_8") // map head

	//mapType := e.getGtype()
	//mapKeyType := mapType.mapKey

	for i, element := range e.elements {
		var offsetKey int = i * mapUnitSize
		var offsetValue int = i*mapUnitSize + 8

		// save key
		emitSaveMapKey(element.key)
		emit("mov %%rax, %%rcx") // copy mallocedaddr
		// append key to tail
		emit("POP_8")                             // map head
		emit("movq %%rcx, %d(%%rax) #", offsetKey) // save key address
		emit("PUSH_8")                            // map head

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

		emit("popq %%rbx") // map head
		emit("movq %%rax, %d(%%rbx) #", offsetValue)
		emit("pushq %%rbx")
	}

	emitCallMalloc(16)
	emit("popq %%rbx") // address (head of the heap)
	emit("movq %%rbx, (%%rax)")
	emit("popq %%rcx") // map len
	emit("movq %%rcx, 8(%%rax)")
}
