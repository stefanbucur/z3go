package main

import (
	"errors"
	"fmt"

	"github.com/stefanbucur/z3go/z3"
)

// Adapted from: http://stackoverflow.com/questions/23451388/z3-sudoku-solver
func solveSudoku(game [][]int) (err error) {
	var conds []*z3.Expr
	vars := make([][]*z3.Expr, 9)

	config := z3.NewConfig()
	ctx := z3.NewContext(config)

	for r := 0; r < 9; r++ {
		vars[r] = make([]*z3.Expr, 9)
		for c := 0; c < 9; c++ {
			vars[r][c] = ctx.IntConst(fmt.Sprintf("x_%d_%d", r+1, c+1))
		}
	}

	// Each cell contains a value in {1, ..., 9}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			conds = append(conds, z3.And(
				z3.Le(ctx.IntVal(1), vars[r][c]),
				z3.Le(vars[r][c], ctx.IntVal(9))))
		}
	}

	// Each row contains a digit at most once
	for r := 0; r < 9; r++ {
		row := make([]*z3.Expr, 9)
		for c := 0; c < 9; c++ {
			row[c] = vars[r][c]
		}
		conds = append(conds, z3.Distinct(row...))
	}

	// Each column contains a digit at most once
	for c := 0; c < 9; c++ {
		column := make([]*z3.Expr, 9)
		for r := 0; r < 9; r++ {
			column[r] = vars[r][c]
		}
		conds = append(conds, z3.Distinct(column...))
	}

	// Each 3x3 square contains a digit at most once
	for sr := 0; sr < 3; sr++ {
		for sc := 0; sc < 3; sc++ {
			square := make([]*z3.Expr, 9)
			for r := 0; r < 3; r++ {
				for c := 0; c < 3; c++ {
					square[3*r+c] = vars[3*sr+r][3*sc+c]
				}
			}
			conds = append(conds, z3.Distinct(square...))
		}
	}

	// Bind the known values
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if game[r][c] != 0 {
				conds = append(conds, z3.Eq(vars[r][c], ctx.IntVal(game[r][c])))
			}
		}
	}

	solver := z3.NewSolver(ctx)
	err = solver.Add(conds...)
	if err != nil {
		return
	}
	result, err := solver.Check()
	if err != nil {
		return
	}
	if result != z3.LTrue {
		err = errors.New("Invalid Sudoku game")
		return
	}
	model := solver.GetModel()
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			fmt.Print(model.Eval(vars[r][c], false), " ")
		}
		fmt.Println()
	}
	return
}
