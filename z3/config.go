package z3

// #cgo LDFLAGS: -lz3
// #include <stdlib.h>
// #include <z3.h>
import "C"
import (
	"runtime"
	"strconv"
	"unsafe"
)

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
