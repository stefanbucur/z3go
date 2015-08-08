package main

import (
	"fmt"

	"github.com/stefanbucur/z3go/z3"
)

/*
# 9x9 matrix of integer variables
X = [ [ Int("x_%s_%s" % (i+1, j+1)) for j in range(9) ]
  for i in range(9) ]

# each cell contains a value in {1, ..., 9}
cells_c  = [ And(1 <= X[i][j], X[i][j] <= 9)
         for i in range(9) for j in range(9) ]

# each row contains a digit at most once
rows_c   = [ Distinct(X[i]) for i in range(9) ]

# each column contains a digit at most once
cols_c   = [ Distinct([ X[i][j] for i in range(9) ])
         for j in range(9) ]

# each 3x3 square contains a digit at most once
sq_c     = [ Distinct([ X[3*i0 + i][3*j0 + j]
                    for i in range(3) for j in range(3) ])
         for i0 in range(3) for j0 in range(3) ]

sudoku_c = cells_c + rows_c + cols_c + sq_c

# sudoku instance, we use '0' for empty cells
instance = ((5,3,0,0,7,0,0,0,0),
        (6,0,0,1,9,5,0,0,0),
        (0,9,8,0,0,0,0,6,0),
        (8,0,0,0,6,0,0,0,3),
        (4,0,0,8,0,3,0,0,1),
        (7,0,0,0,2,0,0,0,6),
        (0,6,0,0,0,0,2,8,0),
        (0,0,0,4,1,9,0,0,5),
        (0,0,0,0,8,0,0,7,9))

instance_c = [ If(instance[i][j] == 0,
              True,
              X[i][j] == instance[i][j])
           for i in range(9) for j in range(9) ]

s = Solver()
s.add(sudoku_c + instance_c)
if s.check() == sat:
m = s.model()
r = [ [ m.evaluate(X[i][j]) for j in range(9) ]
      for i in range(9) ]
print_matrix(r)
else:
print "failed to solve"
*/

// Adapted from: http://stackoverflow.com/questions/23451388/z3-sudoku-solver
func solveSudoku(game [][]int) {
	var vars [9][9]*z3.Expr
	var conds []*z3.Expr

	config := z3.NewConfig()
	ctx := z3.NewContext(config)

	for r := 0; r < 9; r++ {
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
			set := make([]*z3.Expr, 9)
			for r := 0; r < 3; r++ {
				for c := 0; c < 3; c++ {
					set[3*r+c] = vars[3*sr+r][3*sc+c]
				}
			}
			conds = append(conds, z3.Distinct(set...))
		}
	}

	// solver := z3.NewSolver(ctx)
	//solver.add(ctx.BoolVal(true))
	return
}
