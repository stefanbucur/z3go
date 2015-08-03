package z3

// #include <z3.h>
import "C"

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

type AST struct {
	z3val   C.Z3_ast
	context *Context
}

func (ast *AST) ASTKind() (kind ASTKind, err error) {
	kind = ASTKind(C.Z3_get_ast_kind(ast.context.z3val, ast.z3val))
	err = getError(ast.context)
	return
}

func (ast *AST) String() string {
	return C.GoString(C.Z3_ast_to_string(ast.context.z3val, ast.z3val))
}

func (ast *AST) initialize() {
	C.Z3_inc_ref(ast.context.z3val, ast.z3val)
	// TODO: Add a finalizer
}

func (ast *AST) finalize() {
	C.Z3_dec_ref(ast.context.z3val, ast.z3val)
}
