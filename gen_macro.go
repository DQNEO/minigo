package main

func emitMacroDefinitions() {
	emitWithoutIndent("// MACROS")

	macroStart(S("FUNC_PROLOGUE"), S(""))
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")
	macroEnd()

	var i int
	var regi gostring
	for i, regi = range RegsForArguments {
		macroName := Sprintf(S("POP_TO_ARG_%d"), i)
		macroStart(macroName, S(""))
		emit("pop %%%s", regi)
		macroEnd()
	}

	for i, regi = range RegsForArguments {
		macroName := Sprintf(S("PUSH_ARG_%d"), i)
		macroStart(macroName, S(""))
		emit("push %%%s", regi)
		macroEnd()
	}

	macroStart(S("PUSH_8"), S(""))
	emit("push %%rax # primitive")
	macroEnd()

	macroStart(S("PUSH_24"), S(""))
	emit("push %%rax # 1st")
	emit("push %%rbx # 2nd")
	emit("push %%rcx # 3rd")
	macroEnd()

	macroStart(S("PUSH_SLICE"), S(""))
	emit("push %%rax # slice.ptr")
	emit("push %%rbx # slice.len")
	emit("push %%rcx # slice.cap")
	macroEnd()
	macroStart(S("PUSH_INTERFACE"), S(""))
	emit("push %%rax # ifc.1st")
	emit("push %%rbx # ifc.2nd")
	emit("push %%rcx # ifc.3rd")
	macroEnd()

	macroStart(S("POP_8"), S(""))
	emit("pop %%rax # primitive")
	macroEnd()

	macroStart(S("POP_24"), S(""))
	emit("pop %%rcx # 3rd")
	emit("pop %%rbx # 2nd")
	emit("pop %%rax # 1st")
	macroEnd()

	macroStart(S("POP_SLICE"), S(""))
	emit("pop %%rcx # slice.cap")
	emit("pop %%rbx # slice.len")
	emit("pop %%rax # slice.ptr")
	macroEnd()

	macroStart(S("POP_MAP"), S(""))
	emit("pop %%rcx # map.cap")
	emit("pop %%rbx # map.len")
	emit("pop %%rax # map.ptr")
	macroEnd()

	macroStart(S("POP_INTERFACE"), S(""))
	emit("pop %%rcx # ifc.3rd")
	emit("pop %%rbx # ifc.2nd")
	emit("pop %%rax # ifc.1st")
	macroEnd()

	macroStart(S("LOAD_EMPTY_24"), S(""))
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart(S("LOAD_EMPTY_SLICE"), S(""))
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart(S("LOAD_EMPTY_MAP"), S(""))
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart(S("LOAD_EMPTY_INTERFACE"), S(""))
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart(S("LOAD_STRING_LITERAL"), S("slabel"))
	emit("lea \\slabel(%%rip), %%rax")
	macroEnd()

	macroStart(S("LOAD_NUMBER"), S("n"))
	emit("mov $\\n, %%rax")
	macroEnd()

	macroStart(S("STORE_1_TO_LOCAL"), S("offset"))
	emit("mov %%al, \\offset(%%rbp)")
	macroEnd()

	macroStart(S("STORE_8_TO_LOCAL"), S("offset"))
	emit("mov %%rax, \\offset(%%rbp)")
	macroEnd()

	macroStart(S("LOAD_GLOBAL_ADDR"), S("varname, offset"))
	emit("lea \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart(S("LOAD_LOCAL_ADDR"), S("offset"))
	emit("lea \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart(S("LOAD_1_FROM_LOCAL_CAST"), S("offset"))
	emit("movsbq \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart(S("LOAD_1_FROM_LOCAL"), S("offset"))
	emit("mov \\offset(%%rbp), %%al")
	macroEnd()

	macroStart(S("LOAD_8_FROM_LOCAL"), S("offset"))
	emit("mov \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart(S("STORE_1_TO_GLOBAL"), S("varname, offset"))
	emit("mov %%al, \\varname+\\offset(%%rip)")
	macroEnd()

	macroStart(S("STORE_8_TO_GLOBAL"), S("varname, offset"))
	emit("mov %%rax, \\varname+\\offset(%%rip)")
	macroEnd()

	macroStart(S("LOAD_1_FROM_GLOBAL_CAST"), S("varname, offset=0"))
	emit("movsbq \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart(S("LOAD_1_FROM_GLOBAL"), S("varname, offset=0"))
	emit("mov \\varname+\\offset(%%rip), %%al")
	macroEnd()

	macroStart(S("LOAD_8_FROM_GLOBAL"), S("varname, offset=0"))
	emit("mov \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart(S("LOAD_24_BY_DEREF"), S(""))
	emit("mov %d(%%rax), %%rcx", offset16)
	emit("mov %d(%%rax), %%rbx", offset8)
	emit("mov %d(%%rax), %%rax", offset0)
	macroEnd()

	macroStart(S("LOAD_8_BY_DEREF"), S(""))
	emit("mov (%%rax), %%rax")
	macroEnd()

	macroStart(S("LOAD_1_BY_DEREF"), S(""))
	emit("movsbq (%%rax), %%rax")
	macroEnd()

	macroStart(S("LOAD_24_FROM_GLOBAL"), S("varname"))
	emit("mov \\varname+%d(%%rip), %%rax # 1st", offset0)
	emit("mov \\varname+%d(%%rip), %%rbx # 2nd", offset8)
	emit("mov \\varname+%d(%%rip), %%rcx # 3rd", offset16)
	macroEnd()

	macroStart(S("LOAD_24_FROM_LOCAL"), S("offset"))
	emit("mov \\offset+%d(%%rbp), %%rax # 1st", offset0)
	emit("mov \\offset+%d(%%rbp), %%rbx # 2nd", offset8)
	emit("mov \\offset+%d(%%rbp), %%rcx # 3rd", offset16)
	macroEnd()

	macroStart(S("CAST_BYTE_TO_INT"), S(""))
	emit("movzbq %%al, %%rax")
	macroEnd()

	macroStart(S("CMP_EQ_ZERO"), S(""))
	emit("cmp $0, %%rax")
	emit("sete %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart(S("CMP_NE_ZERO"), S(""))
	emit("cmp $0, %%rax")
	emit("setne %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart(S("CMP_FROM_STACK"), S("inst"))
	emit("pop %%rax # right")
	emit("pop %%rcx # left")
	emit("cmp %%rax, %%rcx") // right, left
	emit("\\inst %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart(S("SUM_FROM_STACK"), S(""))
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("add %%rcx , %%rax")
	macroEnd()

	macroStart(S("SUB_FROM_STACK"), S(""))
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("sub %%rcx , %%rax")
	macroEnd()

	macroStart(S("IMUL_FROM_STACK"), S(""))
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("imul %%rcx , %%rax")
	macroEnd()

	macroStart(S("IMUL_NUMBER"), S("n"))
	emit("imul $\\n , %%rax")
	macroEnd()

	macroStart(S("STORE_1_INDIRECT_FROM_STACK"), S(""))
	emit("pop %%rax # where")
	emit("pop %%rcx # what")
	emit("mov %%cl, (%%rax)")
	macroEnd()

	macroStart(S("STORE_8_INDIRECT_FROM_STACK"), S(""))
	emit("pop %%rax # where")
	emit("pop %%rcx # what")
	emit("mov %%rcx, (%%rax)")
	macroEnd()

	macroStart(S("STORE_24_INDIRECT_FROM_STACK"), S(""))
	emit("pop %%rax # target addr")
	emit("pop %%rcx # load RHS value(c)")
	emit("mov %%rcx, 16(%%rax)")
	emit("pop %%rcx # load RHS value(b)")
	emit("mov %%rcx, 8(%%rax)")
	emit("pop %%rcx # load RHS value(a)")
	emit("mov %%rcx, 0(%%rax)")
	macroEnd()

	macroStart(S("ADD_NUMBER"), S("n"))
	emit("add $\\n , %%rax")
	macroEnd()

	macroStart(S("SUB_NUMBER"), S("n"))
	emit("sub $\\n , %%rax")
	macroEnd()

	macroStart(S("FUNCALL"), S("fname"))
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("call \\fname")
	macroEnd()

	macroStart(S("TEST_IT"), S(""))
	emit("test %%rax, %%rax")
	macroEnd()

	macroStart(S("LEAVE_AND_RET"), S(""))
	emit("leave")
	emit("ret")
	macroEnd()
}

func macroStart(name gostring, args gostring) {
	emitWithoutIndent(".macro %s %s", name, args)
}

func macroEnd() {
	emitWithoutIndent(".endm")
	emitNewline()
}
