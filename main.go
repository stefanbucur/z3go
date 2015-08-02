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
	fmt.Println(solver)

	_, err = context.BVSort(32)
	if err != nil {
		fmt.Println("Could not create sort:", err)
		return
	}
}
