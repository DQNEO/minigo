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

	emitWithoutIndent(".macro push_primitive")
	emit("push %%rax # primitive")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro push_slice")
	emit("push %%rax # slice.ptr")
	emit("push %%rbx # slice.len")
	emit("push %%rcx # slice.cap")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro push_map")
	emit("push %%rax # map.ptr")
	emit("push %%rbx # map.len")
	emit("push %%rcx # map.cap")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro push_interface")
	emit("push %%rax # ifc.1st")
	emit("push %%rbx # ifc.2nd")
	emit("push %%rcx # ifc.3rd")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro pop_slice")
	emit("pop %%rcx # slice.cap")
	emit("pop %%rbx # slice.len")
	emit("pop %%rax # slice.ptr")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro pop_map")
	emit("pop %%rcx # map.cap")
	emit("pop %%rbx # map.len")
	emit("pop %%rax # map.ptr")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro pop_interface")
	emit("pop %%rcx # ifc.3rd")
	emit("pop %%rbx # ifc.2nd")
	emit("pop %%rax # ifc.1st")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro load_empty_slice")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro load_empty_map")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro load_empty_interface")
	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")
	emit("mov $0, %%rcx")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro load_string_literal slabel")
	emit("lea \\slabel(%%rip), %%rax")
	emitWithoutIndent(".endm")
	emitNewline()

	emitWithoutIndent(".macro load_number n")
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
}
