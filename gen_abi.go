package main

/**
  Intel® 64 and IA-32 Architectures Software Developer’s Manual
  Combined Volumes: 1, 2A, 2B, 2C, 2D, 3A, 3B, 3C, 3D and 4

  3.4.1.1 General-Purpose Registers in 64-Bit Mode

  In 64-bit mode, there are 16 general purpose registers and the default operand size is 32 bits.
  However, general-purpose registers are able to work with either 32-bit or 64-bit operands.
  If a 32-bit operand size is specified: EAX, EBX, ECX, EDX, EDI, ESI, EBP, ESP, R8D - R15D are available.
  If a 64-bit operand size is specified: RAX, RBX, RCX, RDX, RDI, RSI, RBP, RSP, R8-R15 are available.
  R8D-R15D/R8-R15 represent eight new general-purpose registers.
  All of these registers can be accessed at the byte, word, dword, and qword level.
  REX prefixes are used to generate 64-bit operand sizes or to reference registers R8-R15.
*/

var retRegi [14]string = [14]string{
	"rax",
	"rbx",
	"rcx",
	"rdx",
	"rdi",
	"rsi",
	"r8",
	"r9",
	"r10",
	"r11",
	"r12",
	"r13",
	"r14",
	"r15",
}

var RegsForArguments [12]string = [12]string{
	"rdi",
	"rsi",
	"rdx",
	"rcx",
	"r8",
	"r9",
	"r10",
	"r11",
	"r12",
	"r13",
	"r14",
	"r15",
}

func (f *DeclFunc) prepare() Emitter {

	var params []*ExprVariable

	// prepend receiver to params
	if f.receiver != nil {
		params = []*ExprVariable{f.receiver}
		for _, param := range f.params {
			params = append(params, param)
		}
	} else {
		params = f.params
	}

	var regIndex int
	// offset for params and local variables
	var offset int
	var argRegisters []int
	for _, param := range params {
		var width int
		switch param.getGtype().is24WidthType() {
		case true:
			width = 3
			regIndex += width
			offset -= IntSize * width
			param.offset = offset

			argRegisters = append(argRegisters, regIndex-1)
			argRegisters = append(argRegisters, regIndex-2)
			argRegisters = append(argRegisters, regIndex-3)
		default:
			width = 1
			regIndex += width
			offset -= IntSize * width
			param.offset = offset

			argRegisters = append(argRegisters, regIndex-width)
		}
	}

	var localarea int
	var i int
	var lvar *ExprVariable
	for i, lvar = range f.localvars {
		if lvar.gtype == nil {
			errorft(lvar.token(), "i=%d %s has nil gtype", i, lvar.varname)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size should  not be zero:%s", lvar.gtype.String())
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		//debugf(S("set offset %d to lvar %s, type=%s"), lvar.offset, lvar.varname, lvar.gtype)
	}

	return &funcPrologueEmitter{
		token:        f.token(),
		symbol:       f.getSymbol(),
		argRegisters: argRegisters,
		localvars:    f.localvars,
		localarea:    localarea,
	}
}

type funcPrologueEmitter struct {
	token        *Token
	symbol       string
	argRegisters []int
	localvars    []*ExprVariable
	localarea    int
}

func (fe *funcPrologueEmitter) emit() {
	setPos(fe.token)
	emitWithoutIndent("%s:", fe.symbol)
	emit("FUNC_PROLOGUE")

	if len(fe.argRegisters) > 0 {
		emit("# set params")
	}

	for _, regi := range fe.argRegisters {
		emit("PUSH_ARG_%d", regi)
	}

	if len(fe.localvars) > 0 {
		//emit("# Allocating stack for localvars len=%d", len(fe.localvars))
		for i := len(fe.localvars) - 1; i >= 0; i-- {
			lvar := fe.localvars[i]
			emit("# offset %d variable \"%s\" %s", lvar.offset, lvar.varname, lvar.gtype.String())
		}
		var localarea int = -fe.localarea
		emit("sub $%d, %%rsp # total stack size", localarea)
	}

	emitNewline()
}

func emitCall(symbol string, receiver Expr, args []Expr, params []*ExprVariable) {
	var numRegs int
	if symbol == "" {
		// interface method call
		receiver.emit()
		emit("LOAD_8_BY_DEREF # dereference: convert an interface value to a concrete value")
		emit("PUSH_8 # receiver")
		numRegs = 1
		emitCallInner(numRegs, args, params)

		emit("POP_8 # funcref")
		emit("call *%%rax")
	} else {
		if receiver != nil {
			// method call of a dynamic type
			receiver.emit()
			if receiver.getGtype().is24WidthType() {
				emit("PUSH_24")
				numRegs = 3
			} else {
				emit("PUSH_8")
				numRegs = 1
			}
		}
		emitCallInner(numRegs, args, params)

		emit("FUNCALL %s", symbol)
	}

	emitNewline()
}

func emitCallInner(numRegs int, args []Expr, params []*ExprVariable) {
	var collectVariadicArgs bool // gather variadic args into a slice
	var variadicArgs []Expr
	var arg Expr
	var argIndex int
	var param *ExprVariable
	for argIndex, arg = range args {
		var fromGtype string
		if arg.getGtype() != nil {
			emit("# get fromGtype")
			fromGtype = arg.getGtype().String()
		}
		emit("# from %s", fromGtype)
		if argIndex < len(params) {
			param = params[argIndex]
			if param.isVariadic {
				if _, ok := arg.(*ExprVaArg); !ok {
					collectVariadicArgs = true
				}
			}
		}
		if collectVariadicArgs {
			variadicArgs = append(variadicArgs, arg)
			continue
		}
		var doConvertToInterface bool

		if param != nil {
			emit("# has a corresponding param")

			var fromGtype *Gtype
			if arg.getGtype() != nil {
				fromGtype = arg.getGtype()
				emit("# fromGtype:%s", fromGtype.String())
			}

			var toGtype *Gtype
			if param.getGtype() != nil {
				toGtype = param.getGtype()
				emit("# toGtype:%s", toGtype.String())
			}

			if toGtype != nil && toGtype.getKind() == G_INTERFACE && fromGtype != nil && fromGtype.getKind() != G_INTERFACE {
				doConvertToInterface = true
			}
		}

		emit("# arg %d, doConvertToInterface=%s, collectVariadicArgs=%s",
			argIndex, bool2string(doConvertToInterface), bool2string(collectVariadicArgs))

		if doConvertToInterface {
			emit("# doConvertToInterface !!!")
			emitConversionToInterface(arg)
		} else {
			arg.emit()
		}

		var width int
		if doConvertToInterface || arg.getGtype().is24WidthType() {
			emit("PUSH_24")
			width = 3
		} else {
			emit("PUSH_8")
			width = 1
		}
		numRegs += width
	}

	// check if callee has a variadic
	// https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	// If f is invoked with no actual arguments for p, the value passed to p is nil.
	if !collectVariadicArgs {
		if argIndex+1 < len(params) {
			_param := params[argIndex+1]
			if _param.isVariadic {
				collectVariadicArgs = true
			}
		}
	}

	if collectVariadicArgs {
		emit("# collectVariadicArgs = true")
		lenArgs := len(variadicArgs)
		if lenArgs == 0 {
			emit("LOAD_EMPTY_SLICE")
			emit("PUSH_SLICE")
		} else {
			// var a []interface{}
			for vargIndex, varg := range variadicArgs {
				emit("# emit variadic arg")
				if vargIndex == 0 {
					emit("# make an empty slice to append")
					emit("LOAD_EMPTY_SLICE")
					emit("PUSH_SLICE")
				}
				// conversion : var ifc = x
				if varg.getGtype().getKind() == G_INTERFACE {
					varg.emit()
				} else {
					emitConversionToInterface(varg)
				}
				emit("PUSH_INTERFACE")
				emit("# calling append24")
				emit("POP_TO_ARG_5 # ifc_c")
				emit("POP_TO_ARG_4 # ifc_b")
				emit("POP_TO_ARG_3 # ifc_a")
				emit("POP_TO_ARG_2 # cap")
				emit("POP_TO_ARG_1 # len")
				emit("POP_TO_ARG_0 # ptr")
				emit("FUNCALL %s", getFuncSymbol(IRuntimePath, "append24"))
				emit("PUSH_SLICE")
			}
		}
		numRegs += 3
	}

	emit("# numRegs=%d", numRegs)

	for i := numRegs - 1; i >= 0; i-- {
		if i >= len(RegsForArguments) {
			errorft(args[0].token(), "too many arguments")
		}
		var j int = i
		emit("POP_TO_ARG_%d", j)
	}
}

func (call *IrStaticCall) emit() {
	emit("# emitCallInner %s", call.symbol)
	emitCall(call.symbol, call.receiver, call.args, call.callee.params)
}

func (stmt *StmtReturn) emit() {
	if len(stmt.exprs) == 0 {
		// return void
		emit("mov $0, %%rax")
		stmt.emitDeferAndReturn()
		return
	}

	if len(stmt.exprs) > 7 {
		TBI(stmt.token(), "too many number of arguments")
	}

	var retRegiIndex int
	if len(stmt.exprs) == 1 {
		expr := stmt.exprs[0]
		rettype := stmt.rettypes[0]
		if rettype.getKind() == G_INTERFACE && expr.getGtype().getKind() != G_INTERFACE {
			if expr.getGtype() == nil {
				emit("LOAD_EMPTY_INTERFACE")
			} else {
				emitConversionToInterface(expr)
			}
		} else {
			expr.emit()
			if expr.getGtype() == nil && stmt.rettypes[0].getKind() == G_SLICE {
				emit("LOAD_EMPTY_SLICE")
			}
		}
		stmt.emitDeferAndReturn()
		return
	}
	for i, rettype := range stmt.rettypes {
		expr := stmt.exprs[i]
		expr.emit()
		//		rettype := stmt.rettypes[i]
		if expr.getGtype() == nil && rettype.getKind() == G_SLICE {
			emit("LOAD_EMPTY_SLICE")
		}
		size := rettype.getSize()
		if size < 8 {
			size = 8
		}
		var num64bit int = size / 8 // @TODO odd size
		for j := 0; j < num64bit; j++ {
			reg := retRegi[num64bit-1-j]
			emit("push %%%s", reg)
			retRegiIndex++
		}
	}
	for i := 0; i < retRegiIndex; i++ {
		reg := retRegi[retRegiIndex-1-i]
		emit("pop %%%s", reg)
	}

	stmt.emitDeferAndReturn()
}
