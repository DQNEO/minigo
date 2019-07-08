package main


func eval() {
	fmtPrintf(S("%d\n"), gint)
	fmtPrintf(S("%c\n"), gchar)
	if gtrue {
		fmtPrintf(S("3\n"))
	}
	if !gfalse {
		fmtPrintf(S("4\n"))
	}

	fmtPrintf(S("%d\n"), gstruct.gint)  // 5
	fmtPrintf(S("%c\n"), gstruct.gchar) // 6
	if gstruct.gtrue {
		fmtPrintf(S("7\n"))
	}
	if !gstruct.gfalse {
		fmtPrintf(S("8\n"))
	}

	fmtPrintf(S("%d\n"), garray[0]) // 9
	fmtPrintf(S("%d\n"), garray[1]) // 10
	fmtPrintf(S("%d\n"), garray[2]) // 11

}

func evalnested() {
	// nested data
	fmtPrintf(S("%d\n"), gstruct.inner.gint)        // 12
	fmtPrintf(S("%d\n"), gstruct.inner.inner.gchar) // A
	if gstruct.inner.inner.gtrue == true {
		fmtPrintf(S("14\n"))
	}
}

var garrayommittied [16]int = [16]int{3}

func evalnestedarary() {
	var i int = gpoints[2].Y
	fmtPrintf(S("%d\n"), i+9) //15

	fmtPrintf(S("%d\n"), len(garrayommittied))   // 16
	fmtPrintf(S("%d\n"), garrayommittied[0]+14)  // 17
	fmtPrintf(S("%d\n"), garrayommittied[15]+18) // 18
}

func assign1() {
	gint = 19
	gchar = 66 // B
	gtrue = false
	gfalse = true
	garray = [3]int{23, 24}

	fmtPrintf(S("%d\n"), gint)  // 19
	fmtPrintf(S("%d\n"), gchar) // B
	if !gtrue {
		fmtPrintf(S("21\n")) //21
	}
	if gfalse {
		fmtPrintf(S("22\n")) // 22
	}

	fmtPrintf(S("%d\n"), garray[0])    // 23
	fmtPrintf(S("%d\n"), garray[1])    // 24
	fmtPrintf(S("%d\n"), garray[2]+25) // 25
	gpoint = Point{
		X: 26,
		Y: 27,
	}

	fmtPrintf(S("%d\n"), gpoint.X) // 26
	fmtPrintf(S("%d\n"), gpoint.Y) // 27

}

func assign2() {
	gstructhasslice = StructHasSlice{}
	fmtPrintf(S("%d\n"), len(gstructhasslice.slice)+28) // 28
}

func assign3() {
	gstructhasarray = StructHasArray{
		array: [2]int{28, 29},
	}
	fmtPrintf(S("%d\n"), gstructhasarray.array[1]) // 29
}

func assign4() {
	gstruct = MyStruct{
		gint:   5,
		gchar:  '6',
		gtrue:  true,
		gfalse: false,
		inner: MyInnerStruct{
			gint:  12,
			gtrue: true,
			inner: MyInnerInnerStruct{
				gtrue: true,
				gchar: 67,
			},
		},
	}
	fmtPrintf(S("%d\n"), gstruct.inner.inner.gchar) // C
}

func assign5() {
	gpoints = [...]Point{
		Point{
			X: 26,
			Y: 27,
		},
		Point{
			X: 28,
			Y: 29,
		},
		Point{
			X: 30,
			Y: 31,
		},
	}

	fmtPrintf(S("%d\n"), gpoints[2].Y)   // 31
	fmtPrintf(S("%d\n"), gpoints[1].X+4) // 32
}

/*
func assign6() {
	gpolygon = Polygon{
		line2: [...]Point{
			Point{
				X: 1,
				Y: 2,
			},
			Point{
				X: 3,
				Y: 33,
			},
			Point{
				X: 5,
				Y: 6,
			},
		},
	}

	fmtPrintf(S("%d\n"), gpolygon.line2[1].Y) // 33
}

*/

func main() {
	eval()
	evalnested()
	evalnestedarary()
	assign1()
	assign2()
	assign3()
	assign4()
	assign5()
	//assign6()
}

var gint int = 1
var gchar byte = '2'
var gtrue bool = true
var gfalse bool = false

var gstruct = MyStruct{
	gint:   5,
	gchar:  '6',
	gtrue:  true,
	gfalse: false,
	inner: MyInnerStruct{
		gint:  12,
		gtrue: true,
		inner: MyInnerInnerStruct{
			gtrue: true,
			gchar: 65,
		},
	},
}

var garray [3]int = [...]int{9, 10, 11}

var gpoint Point = Point{
	X: 2,
	Y: 4,
}

var gpoints [3]Point = [...]Point{
	Point{
		X: 1,
		Y: 2,
	},
	Point{
		X: 3,
		Y: 4,
	},
	Point{
		X: 5,
		Y: 6,
	},
}

/*
var gpolygon Polygon = Polygon{
	line2: [...]Point{
		Point{
			X: 1,
			Y: 2,
		},
		Point{
			X: 3,
			Y: 4,
		},
		Point{
			X: 5,
			Y: 6,
		},
	},
}
*/

var gstructhasarray StructHasArray

var gstructhasslice StructHasSlice

type MyStruct struct {
	gint   int
	gchar  byte
	gtrue  bool
	gfalse bool
	inner  MyInnerStruct
}

type MyInnerStruct struct {
	gint   int
	gchar  byte
	gtrue  bool
	gfalse bool
	inner  MyInnerInnerStruct
}

type MyInnerInnerStruct struct {
	gint   int
	gchar  byte
	gtrue  bool
	gfalse bool
}

type Polygon struct {
	line1 [3]Point
	line2 [3]Point
}

type StructHasArray struct {
	array [2]int
}

type StructHasSlice struct {
	slice []int
}

type Point struct {
	X int
	Y int
}
