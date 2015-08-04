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
// Errors

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

func getError(ctx *Context) error {
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
// Contexts

// Context encapsulates a Z3 context.
type Context struct {
	z3val     C.Z3_context
	LastError *Error
}

// Symbol encapsulates a Z3 symbol
type Symbol struct {
	z3val C.Z3_symbol
	ctx   *Context
}

func (ctx *Context) finalize() {
	C.Z3_del_context(ctx.z3val)
}

func (ctx *Context) BVSort(size uint) *Sort {
	z3sort, err := C.Z3_mk_bv_sort(ctx.z3val, C.uint(size)), getError(ctx)
	if err != nil {
		return nil
	}
	return newSort(ctx, z3sort)
}

func (ctx *Context) ArraySort(d *Sort, r *Sort) *Sort {
	z3ds, z3rs := C.Z3_sort(unsafe.Pointer(d.z3val)), C.Z3_sort(unsafe.Pointer(r.z3val))
	z3sort, err := C.Z3_mk_array_sort(ctx.z3val, z3ds, z3rs), getError(ctx)
	if err != nil {
		return nil
	}
	return newSort(ctx, z3sort)
}

// NewContext creates a new Z3 context.
func NewContext(config *Config) *Context {
	ctx := &Context{C.Z3_mk_context_rc(config.z3val), nil}
	C.Z3_set_error_handler(ctx.z3val, nil)
	runtime.SetFinalizer(ctx, (*Context).finalize)
	return ctx
}

func NewStringSymbol(ctx *Context, value string) *Symbol {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	z3sym, err := C.Z3_mk_string_symbol(ctx.z3val, cValue), getError(ctx)
	if err != nil {
		return nil
	}
	return &Symbol{z3sym, ctx}
}

func NewIntSymbol(ctx *Context, value int) *Symbol {
	z3sym, err := C.Z3_mk_int_symbol(ctx.z3val, C.int(value)), getError(ctx)
	if err != nil {
		return nil
	}
	return &Symbol{z3sym, ctx}
}

// -----------------------------------------------------------------------------
// ASTs

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

type (
	AST struct {
		z3val C.Z3_ast
		ctx   *Context
	}

	Expr struct {
		AST
	}

	Sort struct {
		AST
	}
)

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

func newSort(ctx *Context, z3sort C.Z3_sort) *Sort {
	z3ast := C.Z3_ast(unsafe.Pointer(z3sort))
	sort := &Sort{AST{z3ast, ctx}}
	sort.initialize()
	return sort
}

func (sort *Sort) SortKind() SortKind {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	return SortKind(C.Z3_get_sort_kind(sort.ctx.z3val, z3sort))
}

func (sort *Sort) BVSize() (size uint, err error) {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3size, err := C.Z3_get_bv_sort_size(sort.ctx.z3val, z3sort), getError(sort.ctx)
	if err != nil {
		return
	}
	size = uint(z3size)
	return
}

func (sort *Sort) ArrayDomain() *Sort {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3DomSort, err := C.Z3_get_array_sort_domain(sort.ctx.z3val, z3sort), getError(sort.ctx)
	if err != nil {
		return nil
	}
	return newSort(sort.ctx, z3DomSort)
}

func (sort *Sort) ArrayRange() *Sort {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3RangeSort, err := C.Z3_get_array_sort_range(sort.ctx.z3val, z3sort), getError(sort.ctx)
	if err != nil {
		return nil
	}
	return newSort(sort.ctx, z3RangeSort)
}
