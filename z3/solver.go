package z3

// #include <stdlib.h>
// #include <z3.h>
import "C"

// Solver encapsulates a Z3 solver instance.
type Solver struct {
	z3val C.Z3_solver
	ctx   *Context
}

func (solver *Solver) String() string {
	return C.GoString(C.Z3_solver_to_string(solver.ctx.z3val, solver.z3val))
}

func (solver *Solver) Reset() error {
	C.Z3_solver_reset(solver.ctx.z3val, solver.z3val)
	return getError(solver.ctx)
}

// NewSolver creates a new Z3 solver.
func NewSolver(ctx *Context) *Solver {
	solver := &Solver{C.Z3_mk_solver(ctx.z3val), ctx}
	C.Z3_solver_inc_ref(ctx.z3val, solver.z3val)
	return solver
}

// NewSolverForLogic creates a new Z3 solver for a given logic.
func NewSolverForLogic(ctx *Context, logic string) *Solver {
	sym := NewStringSymbol(ctx, logic)
	solver := &Solver{C.Z3_mk_solver_for_logic(ctx.z3val, sym.z3val), ctx}
	C.Z3_solver_inc_ref(ctx.z3val, solver.z3val)
	return solver
}
