package main

import (
	"fmt"

	"github.com/stefanbucur/z3go/z3"
)

func main() {
	config := z3.NewConfig()
	config.SetParamInt("timeout", 1000)

	context := z3.NewContext(config)
	solver, err := z3.NewSolverForLogic(context, "QF_ABV")
	if err != nil {
		fmt.Println("Could not create solver:", err)
		return
	}
	fmt.Println("Created solver", solver)

	sort, err := context.BVSort(32)
	if err != nil {
		fmt.Println("Could not create sort:", err)
		return
	}

	astKind, err := sort.ASTKind()
	sortKind, err := sort.SortKind()
	fmt.Println("Created sort", sort, "of AST kind", astKind, "and sort kind", sortKind)

	_, err = sort.ArrayDomain()
	if err != nil {
		fmt.Println("Expected:", err)
	}
	ssize, err := sort.BVSize()
	fmt.Println("Sort size =", ssize)
}
