package main

import "fmt"

func emitMacroDefinitions() {
	emitWithoutIndent("// MACROS")

	macroStart("FUNC_PROLOGUE", "")
	emit2("push %%rbp")
	emit2("mov %%rsp, %%rbp")
	macroEnd()

	for i, regi := range RegsForArguments {
		macroStart(fmt.Sprintf("POP_TO_ARG_%d", i), "")
		emit2("pop %%%s", gostring(regi))
		macroEnd()
	}

	for i, regi := range RegsForArguments {
		macroStart(fmt.Sprintf("PUSH_ARG_%d", i), "")
		emit2("push %%%s", gostring(regi))
		macroEnd()
	}

	macroStart("PUSH_8", "")
	emit2("push %%rax # primitive")
	macroEnd()

	macroStart("PUSH_24", "")
	emit2("push %%rax # 1st")
	emit2("push %%rbx # 2nd")
	emit2("push %%rcx # 3rd")
	macroEnd()

	macroStart("PUSH_SLICE", "")
	emit2("push %%rax # slice.ptr")
	emit2("push %%rbx # slice.len")
	emit2("push %%rcx # slice.cap")
	macroEnd()
	macroStart("PUSH_INTERFACE", "")
	emit2("push %%rax # ifc.1st")
	emit2("push %%rbx # ifc.2nd")
	emit2("push %%rcx # ifc.3rd")
	macroEnd()

	macroStart("POP_8", "")
	emit2("pop %%rax # primitive")
	macroEnd()

	macroStart("POP_24", "")
	emit2("pop %%rcx # 3rd")
	emit2("pop %%rbx # 2nd")
	emit2("pop %%rax # 1st")
	macroEnd()

	macroStart("POP_SLICE", "")
	emit2("pop %%rcx # slice.cap")
	emit2("pop %%rbx # slice.len")
	emit2("pop %%rax # slice.ptr")
	macroEnd()

	macroStart("POP_MAP", "")
	emit2("pop %%rcx # map.cap")
	emit2("pop %%rbx # map.len")
	emit2("pop %%rax # map.ptr")
	macroEnd()

	macroStart("POP_INTERFACE", "")
	emit2("pop %%rcx # ifc.3rd")
	emit2("pop %%rbx # ifc.2nd")
	emit2("pop %%rax # ifc.1st")
	macroEnd()

	macroStart("LOAD_EMPTY_24", "")
	emit2("mov $0, %%rax")
	emit2("mov $0, %%rbx")
	emit2("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_EMPTY_SLICE", "")
	emit2("mov $0, %%rax")
	emit2("mov $0, %%rbx")
	emit2("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_EMPTY_MAP", "")
	emit2("mov $0, %%rax")
	emit2("mov $0, %%rbx")
	emit2("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_EMPTY_INTERFACE", "")
	emit2("mov $0, %%rax")
	emit2("mov $0, %%rbx")
	emit2("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_STRING_LITERAL", "slabel")
	emit2("lea \\slabel(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_NUMBER", "n")
	emit2("mov $\\n, %%rax")
	macroEnd()

	macroStart("STORE_1_TO_LOCAL", "offset")
	emit2("mov %%al, \\offset(%%rbp)")
	macroEnd()

	macroStart("STORE_8_TO_LOCAL", "offset")
	emit2("mov %%rax, \\offset(%%rbp)")
	macroEnd()

	macroStart("LOAD_GLOBAL_ADDR", "varname, offset")
	emit2("lea \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_LOCAL_ADDR", "offset")
	emit2("lea \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_LOCAL_CAST", "offset")
	emit2("movsbq \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_LOCAL", "offset")
	emit2("mov \\offset(%%rbp), %%al")
	macroEnd()

	macroStart("LOAD_8_FROM_LOCAL", "offset")
	emit2("mov \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("STORE_1_TO_GLOBAL", "varname, offset")
	emit2("mov %%al, \\varname+\\offset(%%rip)")
	macroEnd()

	macroStart("STORE_8_TO_GLOBAL", "varname, offset")
	emit2("mov %%rax, \\varname+\\offset(%%rip)")
	macroEnd()

	macroStart("LOAD_1_FROM_GLOBAL_CAST", "varname, offset=0")
	emit2("movsbq \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_GLOBAL", "varname, offset=0")
	emit2("mov \\varname+\\offset(%%rip), %%al")
	macroEnd()

	macroStart("LOAD_8_FROM_GLOBAL", "varname, offset=0")
	emit2("mov \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_24_BY_DEREF", "")
	emit("mov %d(%%rax), %%rcx", 16)
	emit("mov %d(%%rax), %%rbx", 8)
	emit("mov %d(%%rax), %%rax", 0)
	macroEnd()

	macroStart("LOAD_8_BY_DEREF", "")
	emit("mov (%%rax), %%rax")
	macroEnd()

	macroStart("LOAD_1_BY_DEREF", "")
	emit("movsbq (%%rax), %%rax")
	macroEnd()

	macroStart("LOAD_24_FROM_GLOBAL", "varname")
	emit("mov \\varname+%2d(%%rip), %%rax # 1st", 0)
	emit("mov \\varname+%2d(%%rip), %%rbx # 2nd", 8)
	emit("mov \\varname+%2d(%%rip), %%rcx # 3rd", 16)
	macroEnd()

	macroStart("LOAD_24_FROM_LOCAL", "offset")
	emit("mov \\offset+%2d(%%rbp), %%rax # 1st", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx # 2nd", 8)
	emit("mov \\offset+%2d(%%rbp), %%rcx # 3rd", 16)
	macroEnd()

	macroStart("CAST_BYTE_TO_INT", "")
	emit2("movzbq %%al, %%rax")
	macroEnd()

	macroStart("CMP_EQ_ZERO", "")
	emit2("cmp $0, %%rax")
	emit2("sete %%al")
	emit2("movzb %%al, %%eax")
	macroEnd()

	macroStart("CMP_NE_ZERO", "")
	emit2("cmp $0, %%rax")
	emit2("setne %%al")
	emit2("movzb %%al, %%eax")
	macroEnd()

	macroStart("CMP_FROM_STACK", "inst")
	emit2("pop %%rax # right")
	emit2("pop %%rcx # left")
	emit2("cmp %%rax, %%rcx") // right, left
	emit2("\\inst %%al")
	emit2("movzb %%al, %%eax")
	macroEnd()

	macroStart("SUM_FROM_STACK", "")
	emit2("pop %%rcx")
	emit2("pop %%rax")
	emit2("add %%rcx , %%rax")
	macroEnd()

	macroStart("SUB_FROM_STACK", "")
	emit2("pop %%rcx")
	emit2("pop %%rax")
	emit2("sub %%rcx , %%rax")
	macroEnd()

	macroStart("IMUL_FROM_STACK", "")
	emit2("pop %%rcx")
	emit2("pop %%rax")
	emit2("imul %%rcx , %%rax")
	macroEnd()

	macroStart("IMUL_NUMBER", "n")
	emit2("imul $\\n , %%rax")
	macroEnd()

	macroStart("STORE_1_INDIRECT_FROM_STACK", "")
	emit2("pop %%rax # where")
	emit2("pop %%rcx # what")
	emit2("mov %%cl, (%%rax)")
	macroEnd()

	macroStart("STORE_8_INDIRECT_FROM_STACK", "")
	emit2("pop %%rax # where")
	emit2("pop %%rcx # what")
	emit2("mov %%rcx, (%%rax)")
	macroEnd()

	macroStart("STORE_24_INDIRECT_FROM_STACK", "")
	emit2("pop %%rax # target addr")
	emit2("pop %%rcx # load RHS value(c)")
	emit2("mov %%rcx, 16(%%rax)")
	emit2("pop %%rcx # load RHS value(b)")
	emit2("mov %%rcx, 8(%%rax)")
	emit2("pop %%rcx # load RHS value(a)")
	emit2("mov %%rcx, 0(%%rax)")
	macroEnd()

	macroStart("ADD_NUMBER", "n")
	emit2("add $\\n , %%rax")
	macroEnd()

	macroStart("SUB_NUMBER", "n")
	emit2("sub $\\n , %%rax")
	macroEnd()

	macroStart("FUNCALL", "fname")
	emit2("mov $0, %%rax")
	emit2("mov $0, %%rbx")
	emit2("call \\fname")
	macroEnd()

	macroStart("TEST_IT", "")
	emit2("test %%rax, %%rax")
	macroEnd()

	macroStart("LEAVE_AND_RET", "")
	emit2("leave")
	emit2("ret")
	macroEnd()
}

func macroStart(name string, args string) {
	emitWithoutIndent(".macro %s %s", gostring(name), gostring(args))
}

func macroEnd() {
	emitWithoutIndent(".endm")
	emitNewline()
}
