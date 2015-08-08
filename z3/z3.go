package z3

// #cgo LDFLAGS: -lz3
// #include <stdlib.h>
// #include <z3.h>
import "C"
import (
	"fmt"
	"runtime"
	"strconv"
	"unsafe"
)

// -----------------------------------------------------------------------------
// Configurations

// Config contains Z3 configuration parameters.
type Config struct {
	z3val C.Z3_config
}

// SetParamString sets a configuration parameter using a string value
func (config *Config) SetParamString(id, value string) {
	cID, cValue := C.CString(id), C.CString(value)
	defer func() {
		C.free(unsafe.Pointer(cID))
		C.free(unsafe.Pointer(cValue))
	}()
	C.Z3_set_param_value(config.z3val, cID, cValue)
}

// SetParamInt sets a configuration parameter using an int value
func (config *Config) SetParamInt(id string, value int) {
	config.SetParamString(id, strconv.FormatInt(int64(value), 10))
}

// SetParamBool sets a configuration parameter using a bool value
func (config *Config) SetParamBool(id string, value bool) {
	config.SetParamString(id, strconv.FormatBool(value))
}

func (config *Config) finalize() {
	C.Z3_del_config(config.z3val)
}

// NewConfig creates a new Z3 configuration object.
func NewConfig() *Config {
	config := &Config{C.Z3_mk_config()}
	runtime.SetFinalizer(config, (*Config).finalize)
	return config
}

// -----------------------------------------------------------------------------
// Error codes

// ErrorCode represents a Z3 error code
type ErrorCode int

const (
	OK              ErrorCode = C.Z3_OK                // No error.
	SortError       ErrorCode = C.Z3_SORT_ERROR        // User tried to build an invalid (type incorrect) AST.
	IOB             ErrorCode = C.Z3_IOB               // Index out of bounds
	InvalidArg      ErrorCode = C.Z3_INVALID_ARG       // Invalid argument was provided.
	ParserError     ErrorCode = C.Z3_PARSER_ERROR      // An error occurred when parsing a string or file.
	NoParser        ErrorCode = C.Z3_NO_PARSER         // Parser output is not available, that is, user didn't invoke Z3_parse_smtlib_string or Z3_parse_smtlib_file.
	InvalidPattern  ErrorCode = C.Z3_INVALID_PATTERN   // Invalid pattern was used to build a quantifier.
	MemOutFail      ErrorCode = C.Z3_MEMOUT_FAIL       // A memory allocation failure was encountered.
	FileAccessError ErrorCode = C.Z3_FILE_ACCESS_ERROR // A file could not be accessed.
	InvalidUsage    ErrorCode = C.Z3_INVALID_USAGE     // API call is invalid in the current state.
	InternalFatal   ErrorCode = C.Z3_INTERNAL_FATAL    // An error internal to Z3 occurred.
	DecRefError     ErrorCode = C.Z3_DEC_REF_ERROR     // Trying to decrement the reference counter of an AST that was deleted or the reference counter was not initialized \mlonly.\endmlonly \conly with #Z3_inc_ref.
	Exception       ErrorCode = C.Z3_EXCEPTION         // Internal Z3 exception. Additional details can be retrieved using Z3_get_error_msg.
)

func (code ErrorCode) String() string {
	switch code {
	case OK:
		return "OK"
	case SortError:
		return "SortError"
	case IOB:
		return "IndexOutOfBounds"
	case InvalidArg:
		return "InvalidArg"
	case ParserError:
		return "ParserError"
	case NoParser:
		return "NoParser"
	case InvalidPattern:
		return "InvalidPattern"
	case MemOutFail:
		return "MemOutFail"
	case FileAccessError:
		return "FileAccessError"
	case InvalidUsage:
		return "InvalidUsage"
	case InternalFatal:
		return "InternalFatal"
	case DecRefError:
		return "DecRefError"
	case Exception:
		return "Exception"
	default:
		return "<unknown>"
	}
}

// Error holds the Z3 error code, plus its message
type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// -----------------------------------------------------------------------------
// AST constants

// ASTKind represents a type of AST node
type ASTKind int

// AST type constants
const (
	NumeralAST    ASTKind = C.Z3_NUMERAL_AST
	AppAST        ASTKind = C.Z3_APP_AST
	VarAST        ASTKind = C.Z3_VAR_AST
	QuantifierAST ASTKind = C.Z3_QUANTIFIER_AST
	SortAST       ASTKind = C.Z3_SORT_AST
	FuncDeclAST   ASTKind = C.Z3_FUNC_DECL_AST
	UnknownAST    ASTKind = C.Z3_UNKNOWN_AST
)

func (kind ASTKind) String() string {
	switch kind {
	case NumeralAST:
		return "numeral"
	case AppAST:
		return "app"
	case VarAST:
		return "var"
	case QuantifierAST:
		return "quantifier"
	case SortAST:
		return "sort"
	case FuncDeclAST:
		return "funcdecl"
	case UnknownAST:
		return "unknown"
	default:
		return "<unknown ast>"
	}
}

type SortKind int

const (
	UninterpretedSort SortKind = C.Z3_UNINTERPRETED_SORT
	BoolSort          SortKind = C.Z3_BOOL_SORT
	IntSort           SortKind = C.Z3_INT_SORT
	RealSort          SortKind = C.Z3_REAL_SORT
	BVSort            SortKind = C.Z3_BV_SORT
	ArraySort         SortKind = C.Z3_ARRAY_SORT
	DataTypeSort      SortKind = C.Z3_DATATYPE_SORT
	RelationSort      SortKind = C.Z3_RELATION_SORT
	FiniteDomainSort  SortKind = C.Z3_FINITE_DOMAIN_SORT
	FloatingPointSort SortKind = C.Z3_FLOATING_POINT_SORT
	RoundingModeSort  SortKind = C.Z3_ROUNDING_MODE_SORT
	UnknownSort       SortKind = C.Z3_UNKNOWN_SORT
)

func (sk SortKind) String() string {
	switch sk {
	case UninterpretedSort:
		return "uninterpreted"
	case BoolSort:
		return "bool"
	case IntSort:
		return "int"
	case RealSort:
		return "real"
	case BVSort:
		return "bv"
	case ArraySort:
		return "array"
	case DataTypeSort:
		return "datatype"
	case RelationSort:
		return "relation"
	case FiniteDomainSort:
		return "finitedomain"
	case FloatingPointSort:
		return "floatingpoint"
	case RoundingModeSort:
		return "roundingmode"
	default:
		return "<unknown sort>"
	}
}

// -----------------------------------------------------------------------------
// Contexts

type (
	Context struct {
		z3val     C.Z3_context
		LastError *Error
	}
)

// NewContext creates a new Z3 context.
func NewContext(config *Config) *Context {
	ctx := &Context{C.Z3_mk_context_rc(config.z3val), nil}
	C.Z3_set_error_handler(ctx.z3val, nil)
	runtime.SetFinalizer(ctx, (*Context).finalize)
	return ctx
}

func (ctx *Context) finalize() {
	C.Z3_del_context(ctx.z3val)
}

func (ctx *Context) getError() error {
	ec := ErrorCode(C.Z3_get_error_code(ctx.z3val))
	if ec == OK {
		ctx.LastError = nil
		return nil
	}
	message := C.GoString(C.Z3_get_error_msg_ex(ctx.z3val, C.Z3_error_code(ec)))
	ctx.LastError = &Error{ec, message}
	return ctx.LastError
}

// -----------------------------------------------------------------------------
// Symbols

type Symbol struct {
	z3val C.Z3_symbol
	ctx   *Context
}

func (ctx *Context) NewStringSymbol(value string) *Symbol {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	z3sym, err := C.Z3_mk_string_symbol(ctx.z3val, cValue), ctx.getError()
	if err != nil {
		return nil
	}
	return &Symbol{z3sym, ctx}
}

func (ctx *Context) NewIntSymbol(value int) *Symbol {
	z3sym, err := C.Z3_mk_int_symbol(ctx.z3val, C.int(value)), ctx.getError()
	if err != nil {
		return nil
	}
	return &Symbol{z3sym, ctx}
}

// -----------------------------------------------------------------------------
// ASTs

type AST struct {
	z3val C.Z3_ast
	ctx   *Context
}

func (ast *AST) ASTKind() ASTKind {
	return ASTKind(C.Z3_get_ast_kind(ast.ctx.z3val, ast.z3val))
}

func (ast *AST) String() string {
	return C.GoString(C.Z3_ast_to_string(ast.ctx.z3val, ast.z3val))
}

func (ast *AST) initialize() {
	C.Z3_inc_ref(ast.ctx.z3val, ast.z3val)
	// TODO: Add a finalizer
}

func (ast *AST) finalize() {
	C.Z3_dec_ref(ast.ctx.z3val, ast.z3val)
}

// -----------------------------------------------------------------------------
// Sorts

type Sort struct {
	AST
}

func (sort *Sort) z3sort() C.Z3_sort {
	return C.Z3_sort(unsafe.Pointer(sort.z3val))
}

func (sort *Sort) SortKind() SortKind {
	return SortKind(C.Z3_get_sort_kind(sort.ctx.z3val, sort.z3sort()))
}

func (sort *Sort) BVSize() uint {
	z3size, err := C.Z3_get_bv_sort_size(sort.ctx.z3val, sort.z3sort()), sort.ctx.getError()
	if err != nil {
		return 0
	}
	return uint(z3size)
}

func (sort *Sort) ArrayDomain() *Sort {
	z3DomSort, err := C.Z3_get_array_sort_domain(sort.ctx.z3val, sort.z3sort()), sort.ctx.getError()
	if err != nil {
		return nil
	}
	return sort.ctx.newSort(z3DomSort)
}

func (sort *Sort) ArrayRange() *Sort {
	z3RangeSort, err := C.Z3_get_array_sort_range(sort.ctx.z3val, sort.z3sort()), sort.ctx.getError()
	if err != nil {
		return nil
	}
	return sort.ctx.newSort(z3RangeSort)
}

func (ctx *Context) newSort(z3sort C.Z3_sort) *Sort {
	z3ast := C.Z3_ast(unsafe.Pointer(z3sort))
	sort := &Sort{AST{z3ast, ctx}}
	sort.initialize()
	return sort
}

func (ctx *Context) BoolSort() *Sort {
	z3sort, err := C.Z3_mk_bool_sort(ctx.z3val), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newSort(z3sort)
}

func (ctx *Context) IntSort() *Sort {
	z3sort, err := C.Z3_mk_int_sort(ctx.z3val), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newSort(z3sort)
}

func (ctx *Context) BVSort(size uint) *Sort {
	z3sort, err := C.Z3_mk_bv_sort(ctx.z3val, C.uint(size)), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newSort(z3sort)
}

func (ctx *Context) ArraySort(d *Sort, r *Sort) *Sort {
	z3sort, err := C.Z3_mk_array_sort(ctx.z3val, d.z3sort(), r.z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newSort(z3sort)
}

// -----------------------------------------------------------------------------
// Expressions

type Expr struct {
	AST
}

func (expr *Expr) Sort() *Sort {
	z3sort, err := C.Z3_get_sort(expr.ctx.z3val, expr.z3val), expr.ctx.getError()
	if err != nil {
		return nil
	}
	return expr.ctx.newSort(z3sort)
}

func (ctx *Context) newExpr(z3ast C.Z3_ast) *Expr {
	expr := &Expr{AST{z3ast, ctx}}
	expr.initialize()
	return expr
}

func (ctx *Context) Constant(name string, sort *Sort) *Expr {
	nameSym := ctx.NewStringSymbol(name)
	z3ast, err := C.Z3_mk_const(ctx.z3val, nameSym.z3val, sort.z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

func (ctx *Context) BVConst(name string, size uint) *Expr {
	return ctx.Constant(name, ctx.BVSort(size))
}

func (ctx *Context) BoolConst(name string) *Expr {
	return ctx.Constant(name, ctx.BoolSort())
}

func (ctx *Context) IntConst(name string) *Expr {
	return ctx.Constant(name, ctx.IntSort())
}

func (ctx *Context) BoolVal(b bool) *Expr {
	var z3ast C.Z3_ast
	var err error
	if b {
		z3ast, err = C.Z3_mk_true(ctx.z3val), ctx.getError()
	} else {
		z3ast, err = C.Z3_mk_false(ctx.z3val), ctx.getError()
	}
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

func (ctx *Context) IntVal(n int) *Expr {
	z3ast, err := C.Z3_mk_int(ctx.z3val, C.int(n), ctx.IntSort().z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

func (ctx *Context) UintVal(n uint) *Expr {
	z3ast, err := C.Z3_mk_unsigned_int(ctx.z3val, C.uint(n), ctx.IntSort().z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

func (ctx *Context) BVIntVal(n int, size uint) *Expr {
	z3ast, err := C.Z3_mk_int(ctx.z3val, C.int(n), ctx.BVSort(size).z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

func (ctx *Context) BVUintVal(n uint, size uint) *Expr {
	z3ast, err := C.Z3_mk_unsigned_int(ctx.z3val, C.uint(n), ctx.BVSort(size).z3sort()), ctx.getError()
	if err != nil {
		return nil
	}
	return ctx.newExpr(z3ast)
}

// -----------------------------------------------------------------------------
// Operations

// Array operations

func Store(a, i, v *Expr) *Expr {
	z3ast, err := C.Z3_mk_store(a.ctx.z3val, a.z3val, i.z3val, v.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

// Boolean operators

func Not(a *Expr) *Expr {
	z3ast, err := C.Z3_mk_not(a.ctx.z3val, a.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

func extractASTs(e []*Expr) (asts []C.Z3_ast) {
	asts = make([]C.Z3_ast, len(e))
	for i, expr := range e {
		asts[i] = expr.z3val
	}
	return
}

func And(e ...*Expr) *Expr {
	asts := extractASTs(e)
	z3ast, err := C.Z3_mk_and(e[0].ctx.z3val, C.uint(len(asts)), &asts[0]), e[0].ctx.getError()
	if err != nil {
		return nil
	}
	return e[0].ctx.newExpr(z3ast)
}

func Or(e ...*Expr) *Expr {
	asts := extractASTs(e)
	z3ast, err := C.Z3_mk_or(e[0].ctx.z3val, C.uint(len(asts)), &asts[0]), e[0].ctx.getError()
	if err != nil {
		return nil
	}
	return e[0].ctx.newExpr(z3ast)
}

// Arithmetic operators

// Comparison operators

func Eq(a, b *Expr) *Expr {
	z3ast, err := C.Z3_mk_eq(a.ctx.z3val, a.z3val, b.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

func Distinct(e ...*Expr) *Expr {
	asts := extractASTs(e)
	z3ast, err := C.Z3_mk_distinct(e[0].ctx.z3val, C.uint(len(asts)), &asts[0]), e[0].ctx.getError()
	if err != nil {
		return nil
	}
	return e[0].ctx.newExpr(z3ast)
}

func Lt(a, b *Expr) *Expr {
	z3ast, err := C.Z3_mk_lt(a.ctx.z3val, a.z3val, b.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

func Le(a, b *Expr) *Expr {
	z3ast, err := C.Z3_mk_le(a.ctx.z3val, a.z3val, b.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

func Gt(a, b *Expr) *Expr {
	z3ast, err := C.Z3_mk_gt(a.ctx.z3val, a.z3val, b.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

func Ge(a, b *Expr) *Expr {
	z3ast, err := C.Z3_mk_ge(a.ctx.z3val, a.z3val, b.z3val), a.ctx.getError()
	if err != nil {
		return nil
	}
	return a.ctx.newExpr(z3ast)
}

// ITE

func Ite(c, t, e *Expr) *Expr {
	z3ast, err := C.Z3_mk_ite(c.ctx.z3val, c.z3val, t.z3val, e.z3val), c.ctx.getError()
	if err != nil {
		return nil
	}
	return c.ctx.newExpr(z3ast)
}

// Quantifiers

// -----------------------------------------------------------------------------
// Solvers

type LiftedBool int

const (
	LFalse LiftedBool = C.Z3_L_FALSE
	LUndef LiftedBool = C.Z3_L_UNDEF
	LTrue  LiftedBool = C.Z3_L_TRUE
)

func (lb LiftedBool) String() string {
	switch lb {
	case LFalse:
		return "false"
	case LUndef:
		return "undef"
	case LTrue:
		return "true"
	default:
		return ""
	}
}

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
	return solver.ctx.getError()
}

func (solver *Solver) Push() error {
	C.Z3_solver_push(solver.ctx.z3val, solver.z3val)
	return solver.ctx.getError()
}

func (solver *Solver) Pop(n uint) error {
	C.Z3_solver_pop(solver.ctx.z3val, solver.z3val, C.uint(n))
	return solver.ctx.getError()
}

func (solver *Solver) Check() (result LiftedBool, err error) {
	result = LiftedBool(C.Z3_solver_check(solver.ctx.z3val, solver.z3val))
	err = solver.ctx.getError()
	return
}

func (solver *Solver) Add(a ...*Expr) error {
	for _, expr := range a {
		C.Z3_solver_assert(solver.ctx.z3val, solver.z3val, expr.z3val)
		if err := solver.ctx.getError(); err != nil {
			return err
		}
	}
	return nil
}

// NewSolver creates a new Z3 solver.
func NewSolver(ctx *Context) *Solver {
	solver := &Solver{C.Z3_mk_solver(ctx.z3val), ctx}
	C.Z3_solver_inc_ref(ctx.z3val, solver.z3val)
	return solver
}

// NewSolverForLogic creates a new Z3 solver for a given logic.
func NewSolverForLogic(ctx *Context, logic string) *Solver {
	sym := ctx.NewStringSymbol(logic)
	solver := &Solver{C.Z3_mk_solver_for_logic(ctx.z3val, sym.z3val), ctx}
	C.Z3_solver_inc_ref(ctx.z3val, solver.z3val)
	return solver
}

// -----------------------------------------------------------------------------
// Models

type Model struct {
	z3val C.Z3_model
	ctx   *Context
}

func (ctx *Context) newModel(z3model C.Z3_model) (model *Model) {
	model = &Model{z3model, ctx}
	C.Z3_model_inc_ref(ctx.z3val, z3model)
	return model
}

func (solver *Solver) GetModel() *Model {
	z3model, err := C.Z3_solver_get_model(solver.ctx.z3val, solver.z3val), solver.ctx.getError()
	if err != nil {
		return nil
	}
	return solver.ctx.newModel(z3model)
}

func (model *Model) String() string {
	return C.GoString(C.Z3_model_to_string(model.ctx.z3val, model.z3val))
}

func getZ3Bool(b bool) C.Z3_bool {
	if b {
		return C.Z3_TRUE
	} else {
		return C.Z3_FALSE
	}
}

func (model *Model) Eval(n *Expr, completion bool) (result *Expr) {
	var z3result C.Z3_ast
	status := C.Z3_model_eval(model.ctx.z3val, model.z3val, n.z3val,
		getZ3Bool(completion), &z3result) == C.Z3_TRUE
	if status {
		result = model.ctx.newExpr(z3result)
	}
	return
}
