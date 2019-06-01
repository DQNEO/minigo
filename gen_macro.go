package main

import "fmt"

func emitMacroDefinitions() {
	emitWithoutIndent("// MACROS")

	macroStart("FUNC_PROLOGUE","")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")
	macroEnd()

	for i, regi := range RegsForArguments {
		macroStart(fmt.Sprintf("POP_TO_ARG_%d", i), "")
		emit("pop %%%s", regi)
		macroEnd()
	}

	for i, regi := range RegsForArguments {
		macroStart(fmt.Sprintf("PUSH_ARG_%d", i), "")
		emit("push %%%s", regi)
		macroEnd()
	}

	macroStart("PUSH_8", "")
	emit("push %%rax # primitive")
	macroEnd()

	macroStart("PUSH_24", "")
	emit("push %%rax # 1st")
	emit("push %%rbx # 2nd")
	emit("push %%rcx # 3rd")
	macroEnd()

	macroStart("PUSH_SLICE", "")
	emit("push %%rax # slice.ptr")
	emit("push %%rbx # slice.len")
	emit("push %%rcx # slice.cap")
	macroEnd()

	macroStart("PUSH_MAP", "")
	emit("push %%rax # map.ptr")
	emit("push %%rbx # map.len")
	emit("push %%rcx # map.cap")
	macroEnd()

	macroStart("PUSH_INTERFACE", "")
	emit("push %%rax # ifc.1st")
	emit("push %%rbx # ifc.2nd")
	emit("push %%rcx # ifc.3rd")
	macroEnd()

	macroStart("POP_8", "")
	emit("pop %%rax # primitive")
	macroEnd()

	macroStart("POP_SLICE", "")
	emit("pop %%rcx # slice.cap")
	emit("pop %%rbx # slice.len")
	emit("pop %%rax # slice.ptr")
	macroEnd()

	macroStart("POP_MAP", "")
	emit("pop %%rcx # map.cap")
	emit("pop %%rbx # map.len")
	emit("pop %%rax # map.ptr")
	macroEnd()

	macroStart("POP_INTERFACE", "")
	emit("pop %%rcx # ifc.3rd")
	emit("pop %%rbx # ifc.2nd")
	emit("pop %%rax # ifc.1st")
	macroEnd()

	macroStart("LOAD_EMPTY_SLICE", "")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_EMPTY_MAP", "")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_EMPTY_INTERFACE", "")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	macroEnd()

	macroStart("LOAD_STRING_LITERAL", "slabel")
	emit("lea \\slabel(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_NUMBER",  "n")
	emit("mov $\\n, %%rax")
	macroEnd()

	macroStart("STORE_1_TO_LOCAL", "offset")
	emit("mov %%al, \\offset(%%rbp)")
	macroEnd()

	macroStart("STORE_8_TO_LOCAL", "offset")
	emit("mov %%rax, \\offset(%%rbp)")
	macroEnd()

	macroStart("LOAD_GLOBAL_ADDR", "varname, offset")
	emit("lea \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_LOCAL_ADDR", "offset")
	emit("lea \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_LOCAL_CAST", "offset")
	emit("movsbq \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_LOCAL", "offset")
	emit("mov \\offset(%%rbp), %%al")
	macroEnd()

	macroStart("LOAD_8_FROM_LOCAL", "offset")
	emit("mov \\offset(%%rbp), %%rax")
	macroEnd()

	macroStart("STORE_1_TO_GLOBAL", "varname, offset")
	emit("mov %%al, \\varname+\\offset(%%rip)")
	macroEnd()

	macroStart("STORE_8_TO_GLOBAL", "varname, offset")
	emit("mov %%rax, \\varname+\\offset(%%rip)")
	macroEnd()


	macroStart("LOAD_1_FROM_GLOBAL_CAST", "varname, offset=0")
	emit("movsbq \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_1_FROM_GLOBAL", "varname, offset=0")
	emit("mov \\varname+\\offset(%%rip), %%al")
	macroEnd()

	macroStart("LOAD_8_FROM_GLOBAL", "varname, offset=0")
	emit("mov \\varname+\\offset(%%rip), %%rax")
	macroEnd()

	macroStart("LOAD_24_BY_DEREF", "")
	emit("mov %d(%%rax), %%rcx", 16)
	emit("mov %d(%%rax), %%rbx", 8)
	emit("mov %d(%%rax), %%rax", 0)
	macroEnd()

	macroStart("LOAD_8_BY_DEREF","")
	emit("mov (%%rax), %%rax")
	macroEnd()

	macroStart("LOAD_1_BY_DEREF","")
	emit("movsbq (%%rax), %%rax")
	macroEnd()


	macroStart("LOAD_INTERFACE_FROM_GLOBAL",  "varname")
	emit("mov \\varname+%2d(%%rip), %%rax", 0)
	emit("mov \\varname+%2d(%%rip), %%rbx", ptrSize)
	emit("mov \\varname+%2d(%%rip), %%rcx", ptrSize+ptrSize)
	macroEnd()

	macroStart("LOAD_SLICE_FROM_GLOBAL", "varname")
	emit("mov \\varname+%2d(%%rip), %%rax # ptr", 0)
	emit("mov \\varname+%2d(%%rip), %%rbx # len", ptrSize)
	emit("mov \\varname+%2d(%%rip), %%rcx # cap", ptrSize+IntSize)
	macroEnd()

	macroStart("LOAD_MAP_FROM_GLOBAL", "varname")
	emit("mov \\varname+%2d(%%rip), %%rax # ptr", 0)
	emit("mov \\varname+%2d(%%rip), %%rbx # len", ptrSize)
	emit("mov \\varname+%2d(%%rip), %%rcx # cap", ptrSize+IntSize)
	macroEnd()

	macroStart("LOAD_SLICE_FROM_LOCAL", "offset")
	emit("mov \\offset+%2d(%%rbp), %%rax # ptr", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx # len", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx # cap", ptrSize+IntSize)
	macroEnd()

	macroStart("LOAD_MAP_FROM_LOCAL", "offset")
	emit("mov \\offset+%2d(%%rbp), %%rax # ptr", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx # len", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx # cap", ptrSize+IntSize)
	macroEnd()
	macroStart("LOAD_INTERFACE_FROM_LOCAL",  "offset")
	emit("mov \\offset+%2d(%%rbp), %%rax", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx", ptrSize+ptrSize)
	macroEnd()

	macroStart("CAST_BYTE_TO_INT", "")
	emit("movzbq %%al, %%rax")
	macroEnd()

	macroStart("CMP_EQ_ZERO","")
	emit("cmp $0, %%rax")
	emit("sete %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart("CMP_NE_ZERO","")
	emit("cmp $0, %%rax")
	emit("setne %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart("CMP_FROM_STACK", "inst")
	emit("pop %%rax # right")
	emit("pop %%rcx # left")
	emit("cmp %%rax, %%rcx") // right, left
	emit("\\inst %%al")
	emit("movzb %%al, %%eax")
	macroEnd()

	macroStart("SUM_FROM_STACK", "")
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("add %%rcx , %%rax")
	macroEnd()

	macroStart("SUB_FROM_STACK", "")
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("sub %%rcx , %%rax")
	macroEnd()

	macroStart("IMUL_FROM_STACK", "")
	emit("pop %%rcx")
	emit("pop %%rax")
	emit("imul %%rcx , %%rax")
	macroEnd()

	macroStart("IMUL_NUMBER", "n")
	emit("imul $\\n , %%rax")
	macroEnd()

	macroStart("STORE_1_INDIRECT_FROM_STACK", "")
	emit("pop %%rax")
	emit("pop %%rcx")
	emit("mov %%cl, (%%rax)")
	macroEnd()

	macroStart("STORE_8_INDIRECT_FROM_STACK", "")
	emit("pop %%rax")
	emit("pop %%rcx")
	emit("mov %%rcx, (%%rax)")
	macroEnd()

	macroStart("STORE_24_INDIRECT_FROM_STACK", "")
	emit("pop %%rax # target addr")
	emit("pop %%rcx # load RHS value(c)")
	emit("mov %%rcx, 16(%%rax)")
	emit("pop %%rcx # load RHS value(b)")
	emit("mov %%rcx, 8(%%rax)")
	emit("pop %%rcx # load RHS value(a)")
	emit("mov %%rcx, 0(%%rax)")
	macroEnd()

	macroStart("ADD_NUMBER", "n")
	emit("add $\\n , %%rax")
	macroEnd()

	macroStart("SUB_NUMBER", "n")
	emit("sub $\\n , %%rax")
	macroEnd()

	macroStart("FUNCALL", "fname")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("call \\fname")
	macroEnd()

	macroStart("TEST_IT", "")
	emit("test %%rax, %%rax")
	macroEnd()

	macroStart("LEAVE_AND_RET", "")
	emit("leave")
	emit("ret")
	macroEnd()
}

func macroStart(name string, args string) {
	emitWithoutIndent(".macro %s %s", name, args)
}

func macroEnd() {
	emitWithoutIndent(".endm")
	emitNewline()
}

