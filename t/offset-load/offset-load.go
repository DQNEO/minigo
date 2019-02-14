package main

import "fmt"

type Relation struct {
	gtype *Gtype
}

type Gtype struct {
	typeId int
	relation     *Relation                 // for G_REL
	size         int                       // for scalar type like int, bool, byte, for struct
	origType     *Gtype                    // for pointer
	fields       []*Gtype                  // for struct
	fieldname    string
	offset       int                       // for struct field
	length       int                       // for array, string(len without the terminating \0)
	elementType  *Gtype                    // for array, slice
}


type Ast struct {
	gtype *Gtype
}

func (ast *Ast) getGtype() *Gtype {
	return ast.gtype
}

func f1() int {
	var lhs *Ast = &Ast{
		gtype: &Gtype{
			typeId:11,
			relation:&Relation{
				gtype:&Gtype{

				},
			},
		},
	}

	id := lhs.getGtype().typeId
	return id
}

func main() {
	id := f1()
	fmt.Printf("%d\n", id - 10) // 1
}
