package main

import (
	"strings"
)

// builtin string
var builtinStringKey1 string = "SfmtDumpInterface"
var builtinStringValue1 string = "# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}\\n"
var builtinStringKey2 string = "SfmtDumpSlice"
var builtinStringValue2 string = "# slice = {underlying:%p,len:%d,cap:%d}\\n"

func (program *Program) emitSpecialStrings() {
	// https://sourceware.org/binutils/docs-2.30/as/Data.html#Data
	emit(".data 0")
	emit("# special strings")

	// emit builtin string
	emitWithoutIndent(".%s:", builtinStringKey1)
	emit(".string \"%s\"", builtinStringValue1)
	emitWithoutIndent(".%s:", builtinStringKey2)
	emit(".string \"%s\"", builtinStringValue2)

	// empty string
	eEmptyString.slabel = gostring("empty")
	emitWithoutIndent(".%s:", eEmptyString.slabel)
	emit(".string \"%s\"", eEmptyString.val)
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
	for i := 1; i <= len(program.methodTable); i++ {
		emit(".quad receiverType%d # receiverTypeId:%d", i, i)
	}

	var shortMethodNames []string

	for i := 1; i <= len(program.methodTable); i++ {
		emitWithoutIndent("receiverType%d:", i)
		methods, ok := program.methodTable[i]
		if !ok {
			// This seems not to be harmful? I'm not 100% sure.
			continue
		}
		for _, methodNameFull := range methods {
			if methodNameFull == "." {
				panic("invalid method name")
			}
			splitted := strings.Split(methodNameFull, "$")
			shortMethodName := splitted[1]
			emit(".quad .S.S.%s # key", shortMethodName)
			label := makeLabel()
			gasIndentLevel++
			emit(".data 1")
			emit("%s:", label)
			emit(".quad %s # func addr", methodNameFull)
			gasIndentLevel--
			emit(".data 0")
			emit(".quad %s # func addr addr", label)


			if !in_array(shortMethodName, shortMethodNames) {
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
	emitRuntimeArgs()
	emitMainFunc(program.importOS)
	emitMakeSliceFunc()

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

func emitRuntimeArgs() {
	emitWithoutIndent(".runtime_args:")
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("# set argv, argc, argc")
	emit("mov runtimeArgv(%%rip), %%rax # ptr")
	emit("mov runtimeArgc(%%rip), %%rbx # len")
	emit("mov runtimeArgc(%%rip), %%rcx # cap")

	emitFuncEpilogue(S(".runtime_args_noop_handler"), nil)
}

func emitMainFunc(importOS bool) {
	fname := "main"
	emit(".global	%s", fname)
	emitWithoutIndent("%s:", fname)
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("mov %%rsi, runtimeArgv(%%rip)")
	emit("mov %%rdi, runtimeArgc(%%rip)")
	emit("mov $0, %%rsi")
	emit("mov $0, %%rdi")

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
	emitFuncEpilogue(S("noop_handler"), nil)
}
