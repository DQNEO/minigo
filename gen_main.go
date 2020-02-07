package main

import "fmt"
import "strings"

const __x64_sys_exit int = 60

// builtin string
var builtinStringKey1 string = string("SfmtDumpInterface")
var builtinStringValue1 string = string("# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}")
var builtinStringKey2 string = string("SfmtDumpSlice")
var builtinStringValue2 string = string("# slice = {underlying:%p,len:%d,cap:%d}")

func (program *Program) emitSpecialStrings() {
	// https://sourceware.org/binutils/docs-2.30/as/Data.html#Data
	emit(".data 0")
	emit("# special strings")

	// emit builtin string
	emitWithoutIndent(".%s:", builtinStringKey1)
	emit(".string \"%s\"", builtinStringValue1)
	emitWithoutIndent(".%s:", builtinStringKey2)
	emit(".string \"%s\"", builtinStringValue2)
}

func (program *Program) emitDynamicTypes() {
	emitNewline()
	emit("# Dynamic Types")
	for dynamicTypeId, gs := range symbolTable.uniquedDTypes {
		label := makeDynamicTypeLabel(dynamicTypeId)
		emitWithoutIndent(".%s:", label)
		emit(".string \"%s\"", gs)
	}
}

func (program *Program) emitMethodTable() {
	emitWithoutIndent("#--------------------------------------------------------")
	emit("# Method table")
	emit(".data 0")
	emitWithoutIndent("%s:", "receiverTypes")
	emit(".quad 0 # receiverTypeId:0")
	var maxId int
	var i int
	var id int
	for id, _ = range program.methodTable {
		if maxId < id {
			maxId = id
		}
	}
	for i = 1; i <= maxId; i++ {
		_, ok := program.methodTable[i]
		if ok {
			emit(".quad receiverType%d # receiverTypeId:%d", i, i)
		} else {
			emit(".quad 0")
		}
	}

	var shortMethodNames []string

	for i, v := range program.methodTable {
		emitWithoutIndent("receiverType%d:", i)
		methods := v
		for _, methodNameFull := range methods {
			if methodNameFull == "." {
				panic("invalid method name")
			}
			splitted := strings.Split(methodNameFull, "$")
			var shortMethodName string = splitted[1]
			emit(".quad .S.S.%s # key", shortMethodName)
			label := makeLabel()
			gasIndentLevel++
			emit(".data 1")
			emit("%s:", label)
			emit(".quad %s # func addr", methodNameFull)
			gasIndentLevel--
			emit(".data 0")
			emit(".quad %s # func addr addr", label)

			if !inArray(shortMethodName, shortMethodNames) {
				shortMethodNames = append(shortMethodNames, shortMethodName)
			}
		}
	}

	emitWithoutIndent("#--------------------------------------------------------")
	emitWithoutIndent("# Short method names")
	for _, shortMethodName := range shortMethodNames {
		emit(".data 0")
		emit(".S.S.%s:", shortMethodName)
		gasIndentLevel++
		emit(".data 1")
		emit(".S.%s:", shortMethodName)
		emit(".quad 0") // Any value is ok. This is not referred to.
		gasIndentLevel--
		emit(".data 0")
		emit(".quad .S.%s", shortMethodName)
	}

}

// generate code
func (program *Program) emit() {

	emitMacroDefinitions()

	emit(".data 0")
	program.emitSpecialStrings()
	program.emitDynamicTypes()
	program.emitMethodTable()

	emitWithoutIndent(".text")

	emitMainFunc(program.importOS)
	emitMakeSliceFunc()
	emitSyscallWrapperFunc()

	// emit packages
	for _, pkg := range program.packages {
		emitWithoutIndent("#--------------------------------------------------------")
		emitWithoutIndent("# package %s", pkg.name)
		emitWithoutIndent("# string literals")
		emitWithoutIndent(".data 0")
		for _, ast := range pkg.stringLiterals {
			emitWithoutIndent(".%s:", ast.slabel)
			// https://sourceware.org/binutils/docs-2.30/as/String.html#String
			// the assembler marks the end of each string with a 0 byte.
			emit(".string \"%s\"", ast.val)
		}

		for _, vardecl := range pkg.vars {
			emitNewline()
			vardecl.emit()
		}
		emitNewline()

		emitWithoutIndent(".text")
		for _, funcdecl := range pkg.funcs {
			funcdecl.emit()
			emitNewline()
		}

	}

}

func emitMainFunc(importOS bool) {
	fname := "main"
	emit(".global	%s", fname)
	emitWithoutIndent("%s:", fname)
	emit("mov %%rdi, %%r10")
	emit("mov %%rsi, %%r11")
	emit("mov $0, %%rsi")
	emit("mov $0, %%rdi")

	symbolArgs := fmt.Sprintf("%s.%s", "iruntime", "libcArgs")
	emit("mov %%r11, %s(%%rip)", symbolArgs) // argv
	emit("mov %%r10, %s+8(%%rip)", symbolArgs) // argc
	emit("mov %%r10, %s+16(%%rip)", symbolArgs) // argc

	emit("mov $0, %%r10")
	emit("mov $0, %%r11")

	// init runtime
	emit("# init runtime")
	emit("FUNCALL iruntime.init")

	// init imported packages
	if importOS {
		emit("# init os")
		emit("FUNCALL os.init")
	}

	emitNewline()
	emit("FUNCALL main.main")
	//emit("FUNCALL iruntime.reportMemoryUsage")

	// exit(0)
	emit("mov $%d, %%rax", __x64_sys_exit) // 1st argument
	emit("mov $0,  %%rdi") // int 0
	emit("syscall")
}

func emitSyscallWrapperFunc() {
	// syscall
	emitWithoutIndent("%s:", ".syscallwrapper")
	emit("FUNC_PROLOGUE")
	emitNewline()
	// copied from https://sys.readthedocs.io/en/latest/doc/07_calling_system_calls.html
	emit("movq %%rdi, %%rax") // Syscall number
	emit("movq %%rsi, %%rdi") // shift arg1
	emit("movq %%rdx, %%rsi") // shift arg2
	emit("movq %%rcx, %%rdx") // shift arg3
	emit("movq %%r8, %%r10")  // shift arg4
	emit("movq %%r9, %%r8")   // shift arg5
	emit("movq 8(%%rsp),%%r9")	/* arg6 is on the stack.  */
	emit("syscall")			/* Do the system call.  */
	emit("cmpq $-4095, %%rax")
	emit("LEAVE_AND_RET")			/* Return to caller.  */
	emitNewline()
}
