package z3

/*
#include <z3.h>
*/
import "C"
import "runtime"

// Context encapsulates a Z3 context.
type Context struct {
	z3val C.Z3_context
}

// ContextObject is implemented by all objects that have a Z3 context
type ContextObject interface {
	// GetContext returns the Z3 context of the object
	GetContext() *Context
}

func (context *Context) finalize() {
	C.Z3_del_context(context.z3val)
}

// NewContext creates a new Z3 context.
func NewContext(config *Config) *Context {
	context := &Context{C.Z3_mk_context_rc(config.z3val)}
	runtime.SetFinalizer(context, (*Context).finalize)
	return context
}
