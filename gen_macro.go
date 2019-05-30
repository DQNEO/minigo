package main

func emitMacroDefinitions() {
	emitWithoutIndent("// MACROS")

	emitWithoutIndent(".macro func_prologue")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")
	emitWithoutIndent(".endm")
	emitNewline()

	for i, regi := range RegsForArguments {
		emitWithoutIndent(".macro pop_to_arg_%d", i)
		emitWithoutIndent("pop %%%s", regi)
		emitWithoutIndent(".endm")
		emitNewline()
	}

	for i, regi := range RegsForArguments {
		emitWithoutIndent(".macro push_arg_%d", i)
		emitWithoutIndent("push %%%s", regi)
		emitWithoutIndent(".endm")
		emitNewline()
	}

	emitWithoutIndent(".macro PUSH_PRIMITIVE")
	emit("push %%rax # primitive")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro PUSH_SLICE")
	emit("push %%rax # slice.ptr")
	emit("push %%rbx # slice.len")
	emit("push %%rcx # slice.cap")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro PUSH_MAP")
	emit("push %%rax # map.ptr")
	emit("push %%rbx # map.len")
	emit("push %%rcx # map.cap")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro PUSH_INTERFACE")
	emit("push %%rax # ifc.1st")
	emit("push %%rbx # ifc.2nd")
	emit("push %%rcx # ifc.3rd")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro POP_SLICE")
	emit("pop %%rcx # slice.cap")
	emit("pop %%rbx # slice.len")
	emit("pop %%rax # slice.ptr")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro POP_MAP")
	emit("pop %%rcx # map.cap")
	emit("pop %%rbx # map.len")
	emit("pop %%rax # map.ptr")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro POP_INTERFACE")
	emit("pop %%rcx # ifc.3rd")
	emit("pop %%rbx # ifc.2nd")
	emit("pop %%rax # ifc.1st")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_EMPTY_SLICE")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_EMPTY_MAP")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_EMPTY_INTERFACE")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_STRING_LITERAL slabel")
	emit("lea \\slabel(%%rip), %%rax")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_NUMBER n")
	emit("mov $\\n, %%rax")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro STORE_1_TO_LOCAL offset")
	emit("mov %%al, \\offset(%%rbp)")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro STORE_8_TO_LOCAL offset")
	emit("mov %%rax, \\offset(%%rbp)")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_8_FROM_LOCAL offset")
	emit("mov \\offset(%%rbp), %%rax")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_SLICE_FROM_LOCAL offset")
	emit("mov \\offset+%2d(%%rbp), %%rax # ptr", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx # len", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx # cap", ptrSize+IntSize)
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_MAP_FROM_LOCAL offset")
	emit("mov \\offset+%2d(%%rbp), %%rax # ptr", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx # len", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx # cap", ptrSize+IntSize)
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro LOAD_INTERFACE_FROM_LOCAL offset")
	emit("mov \\offset+%2d(%%rbp), %%rax", 0)
	emit("mov \\offset+%2d(%%rbp), %%rbx", ptrSize)
	emit("mov \\offset+%2d(%%rbp), %%rcx", ptrSize+ptrSize)
	emitWithoutIndent(".endm")
	emitNewline()
}
