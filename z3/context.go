package z3

/*
#include <z3.h>

Z3_sort mk_bv_sort(Z3_context context, unsigned size) {
  Z3_sort result = Z3_mk_bv_sort(context, size);
	Z3_inc_ref(context, (Z3_ast)result);
	return result;
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Context encapsulates a Z3 context.
type Context struct {
	z3val C.Z3_context
}

func (context *Context) finalize() {
	C.Z3_del_context(context.z3val)
}

func (context *Context) BVSort(size uint) (sort *Sort, err error) {
	z3sort, err := C.mk_bv_sort(context.z3val, C.uint(size)), getError(context)
	if err != nil {
		return
	}
	z3ast := C.Z3_ast(unsafe.Pointer(z3sort))
	sort = &Sort{AST{z3ast, context}}
	sort.initialize()
	return
}

// NewContext creates a new Z3 context.
func NewContext(config *Config) *Context {
	context := &Context{C.Z3_mk_context_rc(config.z3val)}
	C.Z3_set_error_handler(context.z3val, nil)
	runtime.SetFinalizer(context, (*Context).finalize)
	return context
}
