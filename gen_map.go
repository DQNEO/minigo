package main

const MAX_METHODS_PER_TYPE int = 128

func (call *IrInterfaceMethodCall) emit() {
	receiver := call.receiver
	var methodName gostring = gostring(call.methodName)
	emit(S("# emit interface method call \"%s\""), methodName)
	mapType := &Gtype{
		kind: G_MAP,
		mapKey: &Gtype{
			kind: G_POINTER,
		},
		mapValue: &Gtype{
			kind: G_POINTER,
		},
	}
	emit(S("# emit receiverTypeId of %s"), receiver.getGtype().String())
	emitOffsetLoad(receiver, ptrSize, ptrSize)
	emit(S("IMUL_NUMBER 8"))
	emit(S("PUSH_8"))

	emit(S("lea receiverTypes(%%rip), %%rax"))
	emit(S("PUSH_8"))
	emit(S("SUM_FROM_STACK"))

	emit(S("# find method %s"), methodName)
	emit(S("mov (%%rax), %%rax")) // address of receiverType
	emit(S("PUSH_8 # map head"))

	emit(S("LOAD_NUMBER %d"), MAX_METHODS_PER_TYPE) // max methods for a type
	emit(S("PUSH_8 # len"))

	emit(S("lea .S.%s, %%rax"), methodName) // index value (addr)
	emit(S("PUSH_8 # map index value"))

	emitMapGet(mapType)

	emit(S("PUSH_8 # funcref"))

	emit(S("mov $0, %%rax"))
	receiverType := receiver.getGtype()
	assert(receiverType.getKind() == G_INTERFACE, nil, S("should be interface"))

	receiver.emit()
	emit(S("LOAD_8_BY_DEREF # dereference: convert an interface value to a concrete value"))

	emit(S("PUSH_8 # receiver"))

	call.emitMethodCall()
}

// emit map index expr
func loadMapIndexExpr(e *ExprIndex) {
	// e.g. x[key]

	_map := e.collection
	// rax: found value (zero if not found)
	// rcx: ok (found: address of the index,  not found:0)
	emit(S("# emit mapData head address"))
	_map.emit()

	// if not nil
	// then emit 24width data
	// else emit 24width 0
	labelNil := makeLabel()
	labelEnd := makeLabel()
	emit(S("TEST_IT # map && map (check if map is nil)"))
	emit(S("je %s # jump if map is nil"), labelNil)
	// not nil case
	emit(S("# not nil"))
	emit(S("LOAD_8_BY_DEREF"))
	emit(S("PUSH_8 # map head"))
	_map.emit()
	emit(S("mov 8(%%rax), %%rax"))
	emit(S("PUSH_8 # len"))
	e.index.emit()
	emit(S("PUSH_8 # index value"))

	emit(S("pop %%rcx # index value"))
	emit(S("pop %%rbx # len"))
	emit(S("pop %%rax # heap"))

	emit(S("jmp %s"), labelEnd)
	// nil case
	emit(S("%s:"), labelNil)
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
	emit(S("%s:"), labelEnd)

	emit(S("PUSH_24"))
	emitMapGet(_map.getGtype())
}

func mapOkRegister(is24Width bool) gostring {
	if is24Width {
		return S("rdx")
	} else {
		return S("rbx")
	}
}

// get map from stack
// r10: map header address")
// r11: map len")
// r12: specified index value")
// r13: loop counter")
func emitMapGet(mapType *Gtype) {

	mapType = mapType.Underlying()
	mapKeyType := mapType.mapKey
	mapValueType := mapType.mapValue
	is24Width := mapValueType.is24WidthType()

	emit(S("# emitMapGet"))

	emit(S("pop %%r12"))
	emit(S("pop %%r11"))
	emit(S("pop %%r10"))

	labelBegin := makeLabel()
	labelEnd := makeLabel()
	labelIncr := makeLabel()

	emit(S("mov $0, %%r13 # init loop counter")) // i = 0

	emit(S("%s: # begin loop "), labelBegin)

	emit(S("push %%r13 # loop counter"))
	emit(S("push %%r11 # map len"))
	emit(S("CMP_FROM_STACK setl"))
	emit(S("TEST_IT"))
	if is24Width {
		emit(S("LOAD_EMPTY_SLICE # NOT FOUND"))
	} else if mapValueType.isString() {
		emitEmptyString()
	} else {
		emit(S("mov $0, %%rax # key not found"))
	}

	okRegister := mapOkRegister(is24Width)
	emit(S("mov $0, %%%s # ok = false"), okRegister)

	emit(S("je %s  # Exit. NOT FOUND IN ALL KEYS."), labelEnd)

	emit(S("# check if key matches"))
	emit(S("mov %%r13, %%rax")) // i
	emit(S("IMUL_NUMBER 16"))   // i * 16
	emit(S("PUSH_8"))

	emit(S("mov %%r10, %%rax")) // head
	emit(S("PUSH_8"))

	emit(S("SUM_FROM_STACK")) // head + i * 16

	emit(S("PUSH_8"))          // index address
	emit(S("LOAD_8_BY_DEREF")) // emit index address

	assert(mapKeyType != nil, nil, S("key kind should not be nil:%s"), mapType.String())

	if mapKeyType.isString() {
		emit(S("push %%r13"))
		emit(S("push %%r11"))
		emit(S("push %%r10"))

		emit(S("LOAD_8_BY_DEREF")) // dereference
		emit(S("PUSH_8"))
		emitConvertCstringFromStackToSlice()
		emit(S("PUSH_SLICE"))

		emit(S("push %%r12"))
		emitConvertCstringFromStackToSlice()
		emit(S("PUSH_SLICE"))

		emitGoStringsEqualFromStack()

		emit(S("pop %%r10"))
		emit(S("pop %%r11"))
		emit(S("pop %%r13"))
	} else {
		emit(S("LOAD_8_BY_DEREF")) // dereference
		// primitive comparison
		emit(S("cmp %%r12, %%rax # compare specifiedvalue vs indexvalue"))
		emit(S("sete %%al"))
		emit(S("movzb %%al, %%eax"))
	}

	emit(S("TEST_IT"))
	emit(S("pop %%rax")) // index address
	emit(S("je %s  # Not match. go to next iteration"), labelIncr)

	emit(S("# Value found!"))
	emit(S("push %%rax # stash key address"))
	emit(S("ADD_NUMBER 8 # value address"))
	emit(S("LOAD_8_BY_DEREF # set the found value address"))
	if mapValueType.is24WidthType() {
		emit(S("LOAD_24_BY_DEREF"))
	} else {
		emit(S("LOAD_8_BY_DEREF"))
	}

	emit(S("mov $1, %%%s # ok = true"), okRegister)
	emit(S("pop %%r12 # key address. will be in map set"))
	emit(S("jmp %s # exit loop"), labelEnd)

	emit(S("%s: # incr"), labelIncr)
	emit(S("add $1, %%r13")) // i++
	emit(S("jmp %s"), labelBegin)

	emit(S("%s: # end loop"), labelEnd)

}

// m[k] = v
// append key and value to the tail of map data, and increment its length
func (e *ExprIndex) emitMapSet(isWidth24 bool) {

	labelAppend := makeLabel()
	labelSave := makeLabel()

	// map get to check if exists
	e.emit()
	// jusdge update or append
	emit(S("cmp $1, %%%s # ok == true"), mapOkRegister(isWidth24))
	emit(S("sete %%al"))
	emit(S("movzb %%al, %%eax"))
	emit(S("TEST_IT"))
	emit(S("je %s  # jump to append if not found"), labelAppend)

	// update
	emit(S("push %%r12")) // push address of the key
	emit(S("jmp %s"), labelSave)

	// append
	emit(S("%s: # append to a map "), labelAppend)
	e.collection.emit() // emit pointer address to %rax
	emit(S("LOAD_8_BY_DEREF"))
	emit(S("PUSH_8"))

	// emit len of the map
	elen := &ExprLen{
		arg: e.collection,
	}
	elen.emit()
	var unitSize int = 2*8
	emit(S("IMUL_NUMBER %d"), unitSize) // distance from head to tail
	emit(S("PUSH_8"))
	emit(S("SUM_FROM_STACK"))
	emit(S("PUSH_8"))

	// map len++
	elen.emit()
	emit(S("ADD_NUMBER 1"))
	emit(S("PUSH_8"))
	e.collection.emit()
	emit(S("pop %%rbx # new len"))
	emit(S("mov %%rbx, 8(%%rax) # update map len"))

	// Save key and value
	emit(S("%s: # end loop"), labelSave)
	e.index.emit()
	emit(S("PUSH_8")) // index value

	// malloc(8)
	emitCallMalloc(8)
	// %%rax : malloced address
	// stack : [map tail address, index value]
	emit(S("pop %%rcx")) // index value

	emit(S("mov %%rcx, (%%rax)")) // save indexvalue to malloced area
	emit(S("mov %%rax, %%rcx"))   // malloced area

	emit(S("POP_8"))              // map tail
	emit(S("mov %%rcx, (%%rax)")) // save indexvalue to map tail
	emit(S("PUSH_8"))             // push map tail

	// save value

	// malloc(8)
	var size int = 8
	if isWidth24 {
		size = 24
	}
	emitCallMalloc(size)

	emit(S("pop %%rcx"))           // map tail address
	emit(S("mov %%rax, 8(%%rcx)")) // set malloced address to tail+8
	emit(S("PUSH_8"))
	if isWidth24 {
		emit(S("STORE_24_INDIRECT_FROM_STACK"))
	} else {
		emit(S("STORE_8_INDIRECT_FROM_STACK"))
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
		op:   gostring("<"),
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

	emit(S("# init index"))
	em.initstmt.emit()

	emit(S("%s: # begin loop "), em.labels.labelBegin)

	em.condition.emit()
	emit(S("TEST_IT"))
	emit(S("je %s  # if false, exit loop"), em.labels.labelEndLoop)

	// set key and value
	em.mapCounter.emit()
	emit(S("IMUL_NUMBER 16"))
	emit(S("PUSH_8 # x"))
	em.rangeexpr.emit() // emit address of map data head
	emit(S("LOAD_8_BY_DEREF"))
	emit(S("PUSH_8 # y"))

	emit(S("SUM_FROM_STACK # x + y"))
	emit(S("LOAD_8_BY_DEREF"))

	emit(S("LOAD_8_BY_DEREF"))
	emitSavePrimitive(em.indexvar)

	if em.valuevar != nil {
		emit(S("# Setting valuevar"))
		emit(S("## rangeexpr.emit()"))
		em.rangeexpr.emit()
		emit(S("LOAD_8_BY_DEREF# map head"))
		emit(S("PUSH_8"))

		emit(S("## mapCounter.emit()"))
		em.mapCounter.emit()
		emit(S("## eval value"))
		emit(S("IMUL_NUMBER 16  # counter * 16"))
		emit(S("ADD_NUMBER 8 # counter * 16 + 8"))
		emit(S("PUSH_8"))

		emit(S("SUM_FROM_STACK"))

		emit(S("LOAD_8_BY_DEREF"))

		if em.valuevar.getGtype().is24WidthType() {
			emit(S("LOAD_24_BY_DEREF"))
			emitSave24(em.valuevar, 0)
		} else {
			emit(S("LOAD_8_BY_DEREF"))
			emitSavePrimitive(em.valuevar)
		}

	}

	em.block.emit()
	emit(S("%s: # end block"), em.labels.labelEndBlock)

	em.indexIncr.emit()

	emit(S("jmp %s"), em.labels.labelBegin)
	emit(S("%s: # end loop"), em.labels.labelEndLoop)
}

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
	emit(S("PUSH_8")) // map head

	//mapType := lit.getGtype()
	//mapKeyType := mapType.mapKey

	for i, element := range lit.elements {
		// alloc key
		if false {
		//	element.key.emit()
		} else {
			element.key.emit()
			emit(S("PUSH_8")) // value of key
			// call malloc for key
			emitCallMalloc(8)
			emit(S("PUSH_8"))

			emit(S("STORE_8_INDIRECT_FROM_STACK")) // save key to heap
		}

		var offsetKey int = i*2*8
		var offsetValue int = i*2*8+8
		emit(S("pop %%rbx"))                         // map head
		emit(S("mov %%rax, %d(%%rbx) #"), offsetKey) // save key address
		emit(S("push %%rbx"))                        // map head

		if element.value.getGtype().getSize() <= 8 {
			element.value.emit()
			emit(S("PUSH_8")) // value of value
			emitCallMalloc(8)
			emit(S("PUSH_8"))
			emit(S("STORE_8_INDIRECT_FROM_STACK")) // save value to heap
		} else if element.value.getGtype().is24WidthType() {
			// rax,rbx,rcx
			element.value.emit()
			emit(S("PUSH_24")) // ptr
			emitCallMalloc(8 * 3)
			emit(S("PUSH_8"))
			emit(S("STORE_24_INDIRECT_FROM_STACK"))
		} else {
			TBI(element.value.token(), S("unable to handle %s"), element.value.getGtype())

		}

		emit(S("pop %%rbx")) // map head
		emit(S("mov %%rax, %d(%%rbx) #"), offsetValue)
		emit(S("push %%rbx"))
	}

	emitCallMalloc(16)
	emit(S("pop %%rbx")) // address (head of the heap)
	emit(S("mov %%rbx, (%%rax)"))
	emit(S("mov $%d, %%rcx"), length) // len
	emit(S("mov %%rcx, 8(%%rax)"))
}
