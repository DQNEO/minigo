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

var retRegi [14]cstring = [14]cstring{
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

var RegsForArguments [12]cstring = [12]cstring{
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
	for _, lvar := range f.localvars {
		if lvar.gtype == nil {
			debugf("%s has nil gtype ", lvar.varname)
		}
		size := lvar.gtype.getSize()
		assert(size != 0, lvar.token(), "size should  not be zero:"+lvar.gtype.String())
		loff := align(size, 8)
		localarea -= loff
		offset -= loff
		lvar.offset = offset
		//debugf("set offset %d to lvar %s, type=%s", lvar.offset, lvar.varname, lvar.gtype)
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
	symbol       gostring
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
			emit("# offset %d variable \"%s\" %s", lvar.offset, gostring(lvar.varname), lvar.gtype.String2())
		}
		var localarea int = -fe.localarea
		emit("sub $%d, %%rsp # total stack size", localarea)
	}

	emitNewline()
}

func (ircall *IrStaticCall) emit() {
	// nothing to do
	emit("# emitCall %s", ircall.symbol)

	var numRegs int
	var param *ExprVariable
	var collectVariadicArgs bool // gather variadic args into a slice
	var variadicArgs []Expr
	var arg Expr
	var argIndex int
	for argIndex, arg = range ircall.args {
		var fromGtype gostring
		if arg.getGtype() != nil {
			emit("# get fromGtype")
			fromGtype = arg.getGtype().String2()
		}
		emit("# from %s", gostring(fromGtype))
		if argIndex < len(ircall.callee.params) {
			param = ircall.callee.params[argIndex]
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

		// do not convert receiver
		if !ircall.isMethodCall || argIndex != 0 {
			if param != nil && !eq(ircall.symbol , "printf") {
				emit("# has a corresponding param")

				var fromGtype *Gtype
				if arg.getGtype() != nil {
					fromGtype = arg.getGtype()
					emit("# fromGtype:%s", fromGtype.String2())
				}

				var toGtype *Gtype
				if param.getGtype() != nil {
					toGtype = param.getGtype()
					emit("# toGtype:%s", toGtype.String2())
				}

				if toGtype != nil && toGtype.getKind() == G_INTERFACE && fromGtype != nil && fromGtype.getKind() != G_INTERFACE {
					doConvertToInterface = true
				}
			}
		}

		if eq(ircall.symbol, ".println") {
			doConvertToInterface = false
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
		if argIndex+1 < len(ircall.callee.params) {
			param = ircall.callee.params[argIndex+1]
			if param.isVariadic {
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
				emit("FUNCALL iruntime.append24")
				emit("PUSH_SLICE")
			}
		}
		numRegs += 3
	}

	for i := numRegs - 1; i >= 0; i-- {
		if i >= len(RegsForArguments) {
			errorft(ircall.args[0].token(), "too many arguments")
		}
		var j int = i
		emit("POP_TO_ARG_%d", j)
	}

	emit("FUNCALL %s", ircall.symbol)
	emitNewline()
}

// @TODO: This is too simple. It should use the same logic as in IrStaticCall for passing args.
func (call *IrInterfaceMethodCall) emitMethodCall() {
	for i, arg := range call.args {
		if _, ok := arg.(*ExprVaArg); ok {
			// skip VaArg for now
			emit("mov $0, %%rax")
		} else {
			arg.emit()
		}
		var no int = i + 2
		emit("PUSH_8 # argument no %d", no)
	}

	var ln int = len(call.args)
	emit("POP_TO_ARG_%d", ln)

	for i := range call.args {
		j := len(call.args) - 1 - i
		var n int = j
		emit("POP_TO_ARG_%d", n)
	}

	emit("POP_8 # funcref")
	emit("call *%%rax")
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
			var reg gostring = gostring(retRegi[num64bit-1-j])
			emit("push %%%s", reg)
			retRegiIndex++
		}
	}
	for i := 0; i < retRegiIndex; i++ {
		var reg gostring = gostring(retRegi[retRegiIndex-1-i])
		emit("pop %%%s", reg)
	}

	stmt.emitDeferAndReturn()
}
