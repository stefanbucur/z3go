package z3

/*
#include <stdlib.h>
#include <z3.h>
*/
import "C"

// Solver encapsulates a Z3 solver instance.
type Solver struct {
	z3val   C.Z3_solver
	context *Context
}

func (solver *Solver) GetContext() *Context {
	return solver.context
}

// NewSolver creates a new Z3 solver.
func NewSolver(context *Context) *Solver {
	solver := &Solver{C.Z3_mk_solver(context.z3val), context}
	C.Z3_solver_inc_ref(context.z3val, solver.z3val)
	return solver
}

// NewSolverForLogic creates a new Z3 solver for a given logic.
func NewSolverForLogic(context *Context, logic string) (*Solver, error) {
	sym, err := NewStringSymbol(context, logic)
	if err != nil {
		return nil, err
	}

	solver := &Solver{C.Z3_mk_solver_for_logic(context.z3val, sym.z3val), context}
	C.Z3_solver_inc_ref(context.z3val, solver.z3val)
	return solver, nil
}
