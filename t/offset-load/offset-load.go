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

func ff1() int {
	var lhs *Ast = &Ast{
		gtype: &Gtype{
			typeId:12,
			relation:&Relation{
				gtype:&Gtype{

				},
			},
		},
	}

	g := lhs.getGtype()
	fields := g.relation.gtype.fields
	fmt.Printf("%d\n", len(fields) + 1) // 0

	for _, fieldtype := range fields {
		fmt.Printf("Error %s\n", fieldtype.fieldname)
	}
	return lhs.getGtype().typeId
}

func f1() {
	id := ff1()
	fmt.Printf("%d\n", id - 10) // 1
}

func main() {
	f1()
}
