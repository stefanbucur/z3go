package z3

// #include <z3.h>
import "C"
import "fmt"

// ErrorCode represents a Z3 error code
type ErrorCode int

const (
	// OK - No error.
	OK ErrorCode = C.Z3_OK

	// SortError - User tried to build an invalid (type incorrect) AST.
	SortError ErrorCode = C.Z3_SORT_ERROR

	// IOB - Index out of bounds
	IOB ErrorCode = C.Z3_IOB

	// InvalidArg - Invalid argument was provided.
	InvalidArg ErrorCode = C.Z3_INVALID_ARG

	// ParserError - An error occurred when parsing a string or file.
	ParserError ErrorCode = C.Z3_PARSER_ERROR

	// NoParser - Parser output is not available, that is, user didn't invoke #Z3_parse_smtlib_string or #Z3_parse_smtlib_file.
	NoParser ErrorCode = C.Z3_NO_PARSER

	// InvalidPattern - Invalid pattern was used to build a quantifier.
	InvalidPattern ErrorCode = C.Z3_INVALID_PATTERN

	// MemOutFail - A memory allocation failure was encountered.
	MemOutFail ErrorCode = C.Z3_MEMOUT_FAIL

	// FileAccessError - A file could not be accessed.
	FileAccessError ErrorCode = C.Z3_FILE_ACCESS_ERROR

	// InvalidUsage - API call is invalid in the current state.
	InvalidUsage ErrorCode = C.Z3_INVALID_USAGE

	// InternalFatal - An error internal to Z3 occurred.
	InternalFatal ErrorCode = C.Z3_INTERNAL_FATAL

	// DecRefError - Trying to decrement the reference counter of an AST that was deleted or the reference counter was not initialized \mlonly.\endmlonly \conly with #Z3_inc_ref.
	DecRefError ErrorCode = C.Z3_DEC_REF_ERROR

	// Exception - Internal Z3 exception. Additional details can be retrieved using #Z3_get_error_msg.
	Exception ErrorCode = C.Z3_EXCEPTION
)

// Error holds the Z3 error code, plus its message
type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func getError(context *Context) *Error {
	ec := ErrorCode(C.Z3_get_error_code(context.z3val))
	if ec == OK {
		return nil
	}
	message := C.GoString(C.Z3_get_error_msg_ex(context.z3val, C.Z3_error_code(ec)))
	return &Error{ec, message}
}
