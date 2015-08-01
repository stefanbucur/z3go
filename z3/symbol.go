package z3

/*
#include <z3.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Symbol encapsulates a Z3 symbol
type Symbol struct {
	z3val   C.Z3_symbol
	context *Context
}

func (symbol *Symbol) GetContext() *Context {
	return symbol.context
}

func NewStringSymbol(context *Context, value string) (*Symbol, error) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	symbol := &Symbol{C.Z3_mk_string_symbol(context.z3val, cValue), context}
	if err := getError(context); err != nil {
		return nil, err
	}
	return symbol, nil
}

func NewIntSymbol(context *Context, value int) (*Symbol, error) {
	symbol := &Symbol{C.Z3_mk_int_symbol(context.z3val, C.int(value)), context}
	if err := getError(context); err != nil {
		return nil, err
	}
	return symbol, nil
}
