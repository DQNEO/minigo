package main

import (
	"./stdlib/fmt"
	"./stdlib/strings"
	"./util"
)

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

			if !util.InArray(shortMethodName, shortMethodNames) {
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

	// emit static function definition
	emitWithoutIndent(".text")
	emitStaticFunctions()

	emit(".data 0")
	program.emitSpecialStrings()
	program.emitDynamicTypes()
	program.emitMethodTable()

	emitWithoutIndent(".text")
	emitMainFunc(program.packages)

	// emit packages
	for _, pkg := range program.packages {
		emitWithoutIndent("#--------------------------------------------------------")
		emitWithoutIndent("# package %s:%s", string(pkg.normalizedPath), pkg.name)
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

func emitMainFunc(packages []*AstPackage) {
	fname := "_start"
	emit(".global	%s", fname)
	emitWithoutIndent("%s:", fname)

	emit("pop %%rax # argc")        // get argc from stack top
	emit("mov %%rsp, %%rbx # argv") // now %rsp value equals to argv

	symbolArgs := fmt.Sprintf("%s.%s", "iruntime", "_argv")
	emit("mov %%rbx, %s(%%rip)", symbolArgs)    // ptr
	emit("mov %%rax, %s+8(%%rip)", symbolArgs)  // len
	emit("mov %%rax, %s+16(%%rip)", symbolArgs) // cap

	emit("mov $0, %%rax")
	emit("mov $0, %%rbx")

	// init imported packages
	for _, pkg := range packages {
		if pkg.hasInit {
			emit("# init %s", string(pkg.normalizedPath))
			emit("FUNCALL %s", getFuncSymbol(pkg.normalizedPath, "init"))
		}
	}

	emitNewline()
	emit("FUNCALL %s", getFuncSymbol("main","main"))
	//emit("FUNCALL iruntime.reportMemoryUsage")

	// exit(0)
	emit("mov $0,  %%rdi")
	emit("FUNCALL %s", getFuncSymbol(IRuntimePath, "exit"))
}
