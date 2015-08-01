package main

import (
	"fmt"

	"github.com/stefanbucur/z3go/z3"
)

func main() {
	config := z3.NewConfig()
	config.SetParamInt("timeout", 1000)

	context := z3.NewContext(config)
	_, err := z3.NewSolverForLogic(context, "QF_ABV")
	if err != nil {
		fmt.Println(err)
	}
}
