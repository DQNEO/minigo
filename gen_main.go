package main

// builtin string
var builtinStringKey1 gostring = gostring("SfmtDumpInterface")
var builtinStringValue1 gostring = gostring("# interface = {ptr:%p,receiverTypeId:%d,dtype:'%s'}")
var builtinStringKey2 gostring = gostring("SfmtDumpSlice")
var builtinStringValue2 gostring = gostring("# slice = {underlying:%p,len:%d,cap:%d}")

func (program *Program) emitSpecialStrings() {
	// https://sourceware.org/binutils/docs-2.30/as/Data.html#Data
	emit(S(".data 0"))
	emit(S("# special strings"))

	// emit builtin string
	emitWithoutIndent(S(".%s:"), builtinStringKey1)
	emit(S(".string \"%s\""), builtinStringValue1)
	emitWithoutIndent(S(".%s:"), builtinStringKey2)
	emit(S(".string \"%s\""), builtinStringValue2)

	// empty string
	eEmptyString.slabel = S("empty")
	emitWithoutIndent(S(".empty:"))
	emit(S(".string \"%s\""), eEmptyString.val)
}

func (program *Program) emitDynamicTypes() {
	emitNewline()
	emit(S("# Dynamic Types"))
	for dynamicTypeId, gs := range symbolTable.uniquedDTypes {
		label := makeDynamicTypeLabel(dynamicTypeId)
		emitWithoutIndent(S(".%s:"), gostring(label))
		emit(S(".string \"%s\""), gostring(gs))
	}
}

func (program *Program) emitMethodTable() {
	emitWithoutIndent(S("#--------------------------------------------------------"))
	emit(S("# Method table"))
	emit(S(".data 0"))
	emitWithoutIndent(S("%s:"), S("receiverTypes"))
	emit(S(".quad 0 # receiverTypeId:0"))
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
			emit(S(".quad receiverType%d # receiverTypeId:%d"), i, i)
		} else {
			emit(S(".quad 0"))
		}
	}

	var shortMethodNames []gostring

	for i,v := range program.methodTable {
		emitWithoutIndent(S("receiverType%d:"), i)
		methods := v
		for _, methodNameFull := range methods {
			if eq(methodNameFull, S(".")) {
				panic(S("invalid method name"))
			}
			splitted := strings_Split(methodNameFull, S("$"))
			var shortMethodName gostring = splitted[1]
			emit(S(".quad .S.S.%s # key"), gostring(shortMethodName))
			label := makeLabel()
			gasIndentLevel++
			emit(S(".data 1"))
			emit(S("%s:"), label)
			emit(S(".quad %s # func addr"), gostring(methodNameFull))
			gasIndentLevel--
			emit(S(".data 0"))
			emit(S(".quad %s # func addr addr"), label)


			if !inArray(shortMethodName, shortMethodNames) {
				shortMethodNames = append(shortMethodNames, shortMethodName)
			}
		}
	}

	emitWithoutIndent(S("#--------------------------------------------------------"))
	emitWithoutIndent(S("# Short method names"))
	for _, shortMethodName := range shortMethodNames {
		emit(S(".data 0"))
		emit(S(".S.S.%s:"), gostring(shortMethodName))
		gasIndentLevel++
		emit(S(".data 1"))
		emit(S(".S.%s:"), gostring(shortMethodName))
		emit(S(".quad 0")) // Any value is ok. This is not referred to.
		gasIndentLevel--
		emit(S(".data 0"))
		emit(S(".quad .S.%s"), gostring(shortMethodName))
	}

}

// generate code
func (program *Program) emit() {

	emitMacroDefinitions()

	emit(S(".data 0"))
	program.emitSpecialStrings()
	program.emitDynamicTypes()
	program.emitMethodTable()

	emitWithoutIndent(S(".text"))

	emitMainFunc(program.importOS)
	emitMakeSliceFunc()

	// emit packages
	for _, pkg := range program.packages {
		emitWithoutIndent(S("#--------------------------------------------------------"))
		emitWithoutIndent(S("# package %s"), gostring(pkg.name))
		emitWithoutIndent(S("# string literals"))
		emitWithoutIndent(S(".data 0"))
		for _, ast := range pkg.stringLiterals {
			emitWithoutIndent(S(".%s:"), ast.slabel)
			// https://sourceware.org/binutils/docs-2.30/as/String.html#String
			// the assembler marks the end of each string with a 0 byte.
			emit(S(".string \"%s\""), ast.val)
		}

		for _, vardecl := range pkg.vars {
			emitNewline()
			vardecl.emit()
		}
		emitNewline()

		emitWithoutIndent(S(".text"))
		for _, funcdecl := range pkg.funcs {
			funcdecl.emit()
			emitNewline()
		}

	}

}

func emitMainFunc(importOS bool) {
	fname := S("main")
	emit(S(".global	%s"), fname)
	emitWithoutIndent(S("%s:"), fname)
	emit(S("push %%rbp"))
	emit(S("mov %%rsp, %%rbp"))

	emit(S("mov %%rsi, libcArgs(%%rip)"))
	emit(S("mov %%rdi, libcArgs+8(%%rip)"))
	emit(S("mov %%rdi, libcArgs+16(%%rip)"))
	emit(S("mov $0, %%rsi"))
	emit(S("mov $0, %%rdi"))

	// init runtime
	emit(S("# init runtime"))
	emit(S("FUNCALL iruntime.init"))

	// init imported packages
	if importOS {
		emit(S("# init os"))
		emit(S("FUNCALL os.init"))
	}

	emitNewline()
	emit(S("FUNCALL main.main"))
	//emit(S("FUNCALL iruntime.reportMemoryUsage"))
	emitFuncEpilogue(S("noop_handler"), nil)
}
