package main

func emitMacroDefinitions() {
	emitWithoutIndent(S("// MACROS"))

	macroStart(S("FUNC_PROLOGUE"), S(""))
	emit(S("push %%rbp"))
	emit(S("mov %%rsp, %%rbp"))
	macroEnd()

	var i int
	var regi gostring
	for i, regi = range RegsForArguments {
		macroName := Sprintf(S("POP_TO_ARG_%d"), i)
		macroStart(macroName, S(""))
		emit(S("pop %%%s"), regi)
		macroEnd()
	}

	for i, regi = range RegsForArguments {
		macroName := Sprintf(S("PUSH_ARG_%d"), i)
		macroStart(macroName, S(""))
		emit(S("push %%%s"), regi)
		macroEnd()
	}

	macroStart(S("PUSH_8"), S(""))
	emit(S("push %%rax # primitive"))
	macroEnd()

	macroStart(S("PUSH_24"), S(""))
	emit(S("push %%rax # 1st"))
	emit(S("push %%rbx # 2nd"))
	emit(S("push %%rcx # 3rd"))
	macroEnd()

	macroStart(S("PUSH_SLICE"), S(""))
	emit(S("push %%rax # slice.ptr"))
	emit(S("push %%rbx # slice.len"))
	emit(S("push %%rcx # slice.cap"))
	macroEnd()
	macroStart(S("PUSH_INTERFACE"), S(""))
	emit(S("push %%rax # ifc.1st"))
	emit(S("push %%rbx # ifc.2nd"))
	emit(S("push %%rcx # ifc.3rd"))
	macroEnd()

	macroStart(S("POP_8"), S(""))
	emit(S("pop %%rax # primitive"))
	macroEnd()

	macroStart(S("POP_24"), S(""))
	emit(S("pop %%rcx # 3rd"))
	emit(S("pop %%rbx # 2nd"))
	emit(S("pop %%rax # 1st"))
	macroEnd()

	macroStart(S("POP_SLICE"), S(""))
	emit(S("pop %%rcx # slice.cap"))
	emit(S("pop %%rbx # slice.len"))
	emit(S("pop %%rax # slice.ptr"))
	macroEnd()

	macroStart(S("POP_MAP"), S(""))
	emit(S("pop %%rcx # map.cap"))
	emit(S("pop %%rbx # map.len"))
	emit(S("pop %%rax # map.ptr"))
	macroEnd()

	macroStart(S("POP_INTERFACE"), S(""))
	emit(S("pop %%rcx # ifc.3rd"))
	emit(S("pop %%rbx # ifc.2nd"))
	emit(S("pop %%rax # ifc.1st"))
	macroEnd()

	macroStart(S("LOAD_EMPTY_24"), S(""))
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
	macroEnd()

	macroStart(S("LOAD_EMPTY_SLICE"), S(""))
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
	macroEnd()

	macroStart(S("LOAD_EMPTY_MAP"), S(""))
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
	macroEnd()

	macroStart(S("LOAD_EMPTY_INTERFACE"), S(""))
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("mov $0, %%rcx"))
	macroEnd()

	macroStart(S("LOAD_STRING_LITERAL"), S("slabel"))
	emit(S("lea \\slabel(%%rip), %%rax"))
	macroEnd()

	macroStart(S("LOAD_NUMBER"), S("n"))
	emit(S("mov $\\n, %%rax"))
	macroEnd()

	macroStart(S("STORE_1_TO_LOCAL"), S("offset"))
	emit(S("mov %%al, \\offset(%%rbp)"))
	macroEnd()

	macroStart(S("STORE_8_TO_LOCAL"), S("offset"))
	emit(S("mov %%rax, \\offset(%%rbp)"))
	macroEnd()

	macroStart(S("LOAD_GLOBAL_ADDR"), S("varname, offset"))
	emit(S("lea \\varname+\\offset(%%rip), %%rax"))
	macroEnd()

	macroStart(S("LOAD_LOCAL_ADDR"), S("offset"))
	emit(S("lea \\offset(%%rbp), %%rax"))
	macroEnd()

	macroStart(S("LOAD_1_FROM_LOCAL_CAST"), S("offset"))
	emit(S("movsbq \\offset(%%rbp), %%rax"))
	macroEnd()

	macroStart(S("LOAD_1_FROM_LOCAL"), S("offset"))
	emit(S("mov \\offset(%%rbp), %%al"))
	macroEnd()

	macroStart(S("LOAD_8_FROM_LOCAL"), S("offset"))
	emit(S("mov \\offset(%%rbp), %%rax"))
	macroEnd()

	macroStart(S("STORE_1_TO_GLOBAL"), S("varname, offset"))
	emit(S("mov %%al, \\varname+\\offset(%%rip)"))
	macroEnd()

	macroStart(S("STORE_8_TO_GLOBAL"), S("varname, offset"))
	emit(S("mov %%rax, \\varname+\\offset(%%rip)"))
	macroEnd()

	macroStart(S("LOAD_1_FROM_GLOBAL_CAST"), S("varname, offset=0"))
	emit(S("movsbq \\varname+\\offset(%%rip), %%rax"))
	macroEnd()

	macroStart(S("LOAD_1_FROM_GLOBAL"), S("varname, offset=0"))
	emit(S("mov \\varname+\\offset(%%rip), %%al"))
	macroEnd()

	macroStart(S("LOAD_8_FROM_GLOBAL"), S("varname, offset=0"))
	emit(S("mov \\varname+\\offset(%%rip), %%rax"))
	macroEnd()

	macroStart(S("LOAD_24_BY_DEREF"), S(""))
	emit(S("mov %d(%%rax), %%rcx"), offset16)
	emit(S("mov %d(%%rax), %%rbx"), offset8)
	emit(S("mov %d(%%rax), %%rax"), offset0)
	macroEnd()

	macroStart(S("LOAD_8_BY_DEREF"), S(""))
	emit(S("mov (%%rax), %%rax"))
	macroEnd()

	macroStart(S("LOAD_1_BY_DEREF"), S(""))
	emit(S("movsbq (%%rax), %%rax"))
	macroEnd()

	macroStart(S("LOAD_24_FROM_GLOBAL"), S("varname"))
	emit(S("mov \\varname+%d(%%rip), %%rax # 1st"), offset0)
	emit(S("mov \\varname+%d(%%rip), %%rbx # 2nd"), offset8)
	emit(S("mov \\varname+%d(%%rip), %%rcx # 3rd"), offset16)
	macroEnd()

	macroStart(S("LOAD_24_FROM_LOCAL"), S("offset"))
	emit(S("mov \\offset+%d(%%rbp), %%rax # 1st"), offset0)
	emit(S("mov \\offset+%d(%%rbp), %%rbx # 2nd"), offset8)
	emit(S("mov \\offset+%d(%%rbp), %%rcx # 3rd"), offset16)
	macroEnd()

	macroStart(S("CAST_BYTE_TO_INT"), S(""))
	emit(S("movzbq %%al, %%rax"))
	macroEnd()

	macroStart(S("CMP_EQ_ZERO"), S(""))
	emit(S("cmp $0, %%rax"))
	emit(S("sete %%al"))
	emit(S("movzb %%al, %%eax"))
	macroEnd()

	macroStart(S("CMP_NE_ZERO"), S(""))
	emit(S("cmp $0, %%rax"))
	emit(S("setne %%al"))
	emit(S("movzb %%al, %%eax"))
	macroEnd()

	macroStart(S("CMP_FROM_STACK"), S("inst"))
	emit(S("pop %%rax # right"))
	emit(S("pop %%rcx # left"))
	emit(S("cmp %%rax, %%rcx")) // right, left
	emit(S("\\inst %%al"))
	emit(S("movzb %%al, %%eax"))
	macroEnd()

	macroStart(S("SUM_FROM_STACK"), S(""))
	emit(S("pop %%rcx"))
	emit(S("pop %%rax"))
	emit(S("add %%rcx , %%rax"))
	macroEnd()

	macroStart(S("SUB_FROM_STACK"), S(""))
	emit(S("pop %%rcx"))
	emit(S("pop %%rax"))
	emit(S("sub %%rcx , %%rax"))
	macroEnd()

	macroStart(S("IMUL_FROM_STACK"), S(""))
	emit(S("pop %%rcx"))
	emit(S("pop %%rax"))
	emit(S("imul %%rcx , %%rax"))
	macroEnd()

	macroStart(S("IMUL_NUMBER"), S("n"))
	emit(S("imul $\\n , %%rax"))
	macroEnd()

	macroStart(S("STORE_1_INDIRECT_FROM_STACK"), S(""))
	emit(S("pop %%rax # where"))
	emit(S("pop %%rcx # what"))
	emit(S("mov %%cl, (%%rax)"))
	macroEnd()

	macroStart(S("STORE_8_INDIRECT_FROM_STACK"), S(""))
	emit(S("pop %%rax # where"))
	emit(S("pop %%rcx # what"))
	emit(S("mov %%rcx, (%%rax)"))
	macroEnd()

	macroStart(S("STORE_24_INDIRECT_FROM_STACK"), S(""))
	emit(S("pop %%rax # target addr"))
	emit(S("pop %%rcx # load RHS value(c)"))
	emit(S("mov %%rcx, 16(%%rax)"))
	emit(S("pop %%rcx # load RHS value(b)"))
	emit(S("mov %%rcx, 8(%%rax)"))
	emit(S("pop %%rcx # load RHS value(a)"))
	emit(S("mov %%rcx, 0(%%rax)"))
	macroEnd()

	macroStart(S("ADD_NUMBER"), S("n"))
	emit(S("add $\\n , %%rax"))
	macroEnd()

	macroStart(S("SUB_NUMBER"), S("n"))
	emit(S("sub $\\n , %%rax"))
	macroEnd()

	macroStart(S("FUNCALL"), S("fname"))
	emit(S("mov $0, %%rax"))
	emit(S("mov $0, %%rbx"))
	emit(S("call \\fname"))
	macroEnd()

	macroStart(S("TEST_IT"), S(""))
	emit(S("test %%rax, %%rax"))
	macroEnd()

	macroStart(S("LEAVE_AND_RET"), S(""))
	emit(S("leave"))
	emit(S("ret"))
	macroEnd()
}

func macroStart(name gostring, args gostring) {
	emitWithoutIndent(S(".macro %s %s"), name, args)
}

func macroEnd() {
	emitWithoutIndent(S(".endm"))
	emitNewline()
}
