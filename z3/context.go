package z3

// #include <z3.h>
import "C"
import "runtime"

// Context encapsulates a Z3 context.
type Context struct {
	z3val C.Z3_context
}

func (context *Context) finalize() {
	C.Z3_del_context(context.z3val)
}

func (context *Context) BVSort(size uint) (sort *Sort, err error) {
	z3sort, err := C.Z3_mk_bv_sort(context.z3val, C.uint(size)), getError(context)
	if err != nil {
		return
	}
	sort = newSort(context, z3sort)
	return
}

// NewContext creates a new Z3 context.
func NewContext(config *Config) *Context {
	context := &Context{C.Z3_mk_context_rc(config.z3val)}
	C.Z3_set_error_handler(context.z3val, nil)
	runtime.SetFinalizer(context, (*Context).finalize)
	return context
}
