package main

import "fmt"

func eval() {
	fmt.Printf("%d\n", gint)
	fmt.Printf("%c\n", gchar)
	if gtrue {
		fmt.Printf("3\n")
	}
	if !gfalse {
		fmt.Printf("4\n")
	}

	fmt.Printf("%d\n", gstruct.gint)  // 5
	fmt.Printf("%c\n", gstruct.gchar) // 6
	if gstruct.gtrue {
		fmt.Printf("7\n")
	}
	if !gstruct.gfalse {
		fmt.Printf("8\n")
	}

	fmt.Printf("%d\n", garray[0]) // 9
	fmt.Printf("%d\n", garray[1]) // 10
	fmt.Printf("%d\n", garray[2]) // 11

}

func evalnested() {
	// nested data
	fmt.Printf("%d\n", gstruct.inner.gint) // 12
	fmt.Printf("%d\n", gstruct.inner.inner.gchar) // 13
	if gstruct.inner.inner.gtrue == true {
		fmt.Printf("14\n")
	}
}

var garrayommittied [16]int = [16]int{3}

func evalnestedarary() {
	var i int = gpoints[2].Y
	fmt.Printf("%d\n", i + 9) //15

	fmt.Printf("%d\n", len(garrayommittied)) // 16
	fmt.Printf("%d\n", garrayommittied[0] + 14) // 17
	fmt.Printf("%d\n", garrayommittied[15] + 18) // 18
}

func assign1() {
	gint = 19
	gchar = 20
	gtrue = false
	gfalse = true
	garray = [3]int{23, 24}

	fmt.Printf("%d\n", gint)  // 19
	fmt.Printf("%d\n", gchar) // 20
	if !gtrue {
		fmt.Printf("21\n") //21
	}
	if gfalse {
		fmt.Printf("22\n") // 22
	}

	fmt.Printf("%d\n", garray[0])    // 23
	fmt.Printf("%d\n", garray[1])    // 24
	fmt.Printf("%d\n", garray[2]+25) // 25
	gpoint = Point{
		X:26,
		Y:27,
	}

	fmt.Printf("%d\n", gpoint.X) // 26
	fmt.Printf("%d\n", gpoint.Y) // 27

}

func assign2() {
	gstructhasslice = StructHasSlice{}
	fmt.Printf("%d\n", len(gstructhasslice.slice) + 28) // 28

	gstructhasarray = StructHasArray{
		array: [2]int{28,29},
	}
	//fmt.Printf("%d\n", gstructhasarray.array[0]) // 28

	/*
	*/
	/*
	gstruct = MyStruct{
		gint:5,
		gchar:'6',
		gtrue:true,
		gfalse:false,
		inner: MyInnerStruct{
			gint:12,
			gtrue: true,
			inner: MyInnerInnerStruct{
				gtrue:true,
				gchar:13,
			},
		},
	}

	gpoints = [...]Point{
		Point{
			X:1,
			Y:2,
		},
		Point{
			X:3,
			Y:4,
		},
		Point{
			X:5,
			Y:6,
		},
	}

	gpolygon = Polygon{
		line2: [...]Point{
			Point{
				X:1,
				Y:2,
			},
			Point{
				X:3,
				Y:4,
			},
			Point{
				X:5,
				Y:6,
			},
		},
	}
	*/
}

func main() {
	eval()
	evalnested()
	evalnestedarary()
	assign1()
	assign2()
}


var gint int = 1
var gchar byte = '2'
var gtrue bool = true
var gfalse bool = false

var gstruct = MyStruct{
	gint:5,
	gchar:'6',
	gtrue:true,
	gfalse:false,
	inner: MyInnerStruct{
		gint:12,
		gtrue: true,
		inner: MyInnerInnerStruct{
			gtrue:true,
			gchar:13,
		},
	},
}

var garray [3]int = [...]int{9,10,11}

var gpoint Point = Point{
	X:2,
	Y:4,
}

var gpoints [3]Point = [...]Point{
	Point{
		X:1,
		Y:2,
	},
	Point{
		X:3,
		Y:4,
	},
	Point{
		X:5,
		Y:6,
	},
}

var gpolygon Polygon = Polygon{
	line2: [...]Point{
		Point{
			X:1,
			Y:2,
		},
		Point{
			X:3,
			Y:4,
		},
		Point{
			X:5,
			Y:6,
		},
	},
}

var gstructhasarray StructHasArray

var gstructhasslice StructHasSlice

type MyStruct struct {
	gint int
	gchar byte
	gtrue bool
	gfalse bool
	inner MyInnerStruct
}

type MyInnerStruct struct {
	gint int
	gchar byte
	gtrue bool
	gfalse bool
	inner MyInnerInnerStruct
}

type MyInnerInnerStruct struct {
	gint int
	gchar byte
	gtrue bool
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
