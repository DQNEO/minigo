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
	emit2(".data 0")
	emit2("# special strings")

	// emit builtin string
	emitWithoutIndent(".%s:", builtinStringKey1)
	emit2(".string \"%s\"", builtinStringValue1)
	emitWithoutIndent(".%s:", builtinStringKey2)
	emit2(".string \"%s\"", builtinStringValue2)

	// empty string
	eEmptyString.slabel = gostring("empty")
	emitWithoutIndent(".%s:", eEmptyString.slabel)
	emit2(".string \"%s\"", eEmptyString.val)
}

func (program *Program) emitDynamicTypes() {
	emitNewline()
	emit2("# Dynamic Types")
	for dynamicTypeId, gs := range symbolTable.uniquedDTypes {
		label := makeDynamicTypeLabel(dynamicTypeId)
		emitWithoutIndent(".%s:", label)
		emit2(".string \"%s\"", gs)
	}
}

func (program *Program) emitMethodTable() {
	emitWithoutIndent("#--------------------------------------------------------")
	emit2("# Method table")
	emit2(".data 0")
	emitWithoutIndent("%s:", "receiverTypes")
	emit2(".quad 0 # receiverTypeId:0")
	var i int
	for i = 1; i <= len(program.methodTable); i++ {
		emit2(".quad receiverType%d # receiverTypeId:%d", i, i)
	}

	var shortMethodNames []string

	for i = 1; i <= len(program.methodTable); i++ {
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
			var shortMethodName = splitted[1]
			emit2(".quad .S.S.%s # key", gostring(shortMethodName))
			label := makeLabel()
			gasIndentLevel++
			emit2(".data 1")
			emit2("%s:", label)
			emit2(".quad %s # func addr", gostring(methodNameFull))
			gasIndentLevel--
			emit2(".data 0")
			emit2(".quad %s # func addr addr", label)


			if !in_array(shortMethodName, shortMethodNames) {
				shortMethodNames = append(shortMethodNames, shortMethodName)
			}
		}
	}

	emitWithoutIndent("#--------------------------------------------------------")
	emitWithoutIndent("# Short method names")
	for _, shortMethodName := range shortMethodNames {
		emit2(".data 0")
		emit2(".S.S.%s:", gostring(shortMethodName))
		gasIndentLevel++
		emit2(".data 1")
		emit2(".S.%s:", gostring(shortMethodName))
		emit2(".quad 0") // Any value is ok. This is not referred to.
		gasIndentLevel--
		emit2(".data 0")
		emit2(".quad .S.%s", gostring(shortMethodName))
	}

}

// generate code
func (program *Program) emit() {

	emitMacroDefinitions()

	emit2(".data 0")
	program.emitSpecialStrings()
	program.emitDynamicTypes()
	program.emitMethodTable()

	emitWithoutIndent2(".text")
	emitRuntimeArgs()
	emitMainFunc(program.importOS)
	emitMakeSliceFunc()

	// emit packages
	for _, pkg := range program.packages {
		emitWithoutIndent2("#--------------------------------------------------------")
		emitWithoutIndent2("# package %s", gostring(pkg.name))
		emitWithoutIndent2("# string literals")
		emitWithoutIndent2(".data 0")
		for _, ast := range pkg.stringLiterals {
			emitWithoutIndent2(".%s:", ast.slabel)
			// https://sourceware.org/binutils/docs-2.30/as/String.html#String
			// the assembler marks the end of each string with a 0 byte.
			emit2(".string \"%s\"", ast.val)
		}

		for _, vardecl := range pkg.vars {
			emitNewline()
			vardecl.emit()
		}
		emitNewline()

		emitWithoutIndent2(".text")
		for _, funcdecl := range pkg.funcs {
			funcdecl.emit()
			emitNewline()
		}

	}

}

func emitRuntimeArgs() {
	emitWithoutIndent(".runtime_args:")
	emit2("push %%rbp")
	emit2("mov %%rsp, %%rbp")

	emit2("# set argv, argc, argc")
	emit2("mov runtimeArgv(%%rip), %%rax # ptr")
	emit2("mov runtimeArgc(%%rip), %%rbx # len")
	emit2("mov runtimeArgc(%%rip), %%rcx # cap")

	emitFuncEpilogue(S(".runtime_args_noop_handler"), nil)
}

func emitMainFunc(importOS bool) {
	fname := S("main")
	emit2(".global	%s", fname)
	emitWithoutIndent("%s:", fname)
	emit2("push %%rbp")
	emit2("mov %%rsp, %%rbp")

	emit2("mov %%rsi, runtimeArgv(%%rip)")
	emit2("mov %%rdi, runtimeArgc(%%rip)")
	emit2("mov $0, %%rsi")
	emit2("mov $0, %%rdi")

	// init runtime
	emit2("# init runtime")
	emit2("FUNCALL iruntime.init")

	// init imported packages
	if importOS {
		emit2("# init os")
		emit2("FUNCALL os.init")
	}

	emitNewline()
	emit2("FUNCALL main.main")
	//emit2("FUNCALL iruntime.reportMemoryUsage")
	emitFuncEpilogue(S("noop_handler"), nil)
}
