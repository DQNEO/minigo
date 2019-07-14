package main

import "strings"

// builtin string
var builtinStringKey1 bytes = bytes("SfmtDumpInterface")
var builtinStringValue1 bytes = bytes("# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}")
var builtinStringKey2 bytes = bytes("SfmtDumpSlice")
var builtinStringValue2 bytes = bytes("# slice = {underlying:%p,len:%d,cap:%d}")

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
		emitWithoutIndent(".%s:", bytes(label))
		emit(".string \"%s\"", bytes(gs))
	}
}

func (program *Program) emitMethodTable() {
	emitWithoutIndent("#--------------------------------------------------------")
	emit("# Method table")
	emit(".data 0")
	emitWithoutIndent("%s:", S("receiverTypes"))
	emit(".quad 0 # receiverTypeId:0")
	var maxId int
	var i int
	var id int
	for id,_ = range program.methodTable {
		if maxId < id {
			maxId = id
		}
	}
	for i=1;i<=maxId;i++ {
		_,ok := program.methodTable[i]
		if ok {
			emit(".quad receiverType%d # receiverTypeId:%d", i, i)
		} else {
			emit(".quad 0")
		}
	}

	var shortMethodNames []bytes

	for i,v := range program.methodTable {
		emitWithoutIndent("receiverType%d:", i)
		methods := v
		for _, methodNameFull := range methods {
			if eq(methodNameFull, S(".")) {
				panic(S("invalid method name"))
			}
			splitted := strings.Split(string(methodNameFull), "$")
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


			if !inArray(bytes(shortMethodName), shortMethodNames) {
				shortMethodNames = append(shortMethodNames, bytes(shortMethodName))
			}
		}
	}

	emitWithoutIndent("#--------------------------------------------------------")
	emitWithoutIndent("# Short method names")
	for _, shortMethodName := range shortMethodNames {
		emit(".data 0")
		emit(".S.S.%s:", bytes(shortMethodName))
		gasIndentLevel++
		emit(".data 1")
		emit(".S.%s:", bytes(shortMethodName))
		emit(".quad 0") // Any value is ok. This is not referred to.
		gasIndentLevel--
		emit(".data 0")
		emit(".quad .S.%s", bytes(shortMethodName))
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

	// emit packages
	for _, pkg := range program.packages {
		emitWithoutIndent("#--------------------------------------------------------")
		emitWithoutIndent("# package %s", bytes(pkg.name))
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
	emit("push %%rbp")
	emit("mov %%rsp, %%rbp")

	emit("mov %%rsi, libcArgs(%%rip)")
	emit("mov %%rdi, libcArgs+8(%%rip)")
	emit("mov %%rdi, libcArgs+16(%%rip)")
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
	emitFuncEpilogue("noop_handler", nil)
}
