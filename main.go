package main

import (
	"fmt"

	"github.com/stefanbucur/z3go/z3"
)

func main() {
	config := z3.NewConfig()
	config.SetParamInt("timeout", 1000)

	context := z3.NewContext(config)
	solver := z3.NewSolverForLogic(context, "QF_ABV")
	fmt.Println("Created solver", solver)

	sort := context.BVSort(32)
	fmt.Printf("Created sort %s of AST kind %q and sort kind %q\n",
		sort, sort.ASTKind(), sort.SortKind())
}
