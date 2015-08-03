package z3

// #include <z3.h>
import "C"
import "unsafe"

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

type Sort struct {
	AST
}

func newSort(context *Context, z3sort C.Z3_sort) *Sort {
	z3ast := C.Z3_ast(unsafe.Pointer(z3sort))
	sort := &Sort{AST{z3ast, context}}
	sort.initialize()
	return sort
}

func (sort *Sort) SortKind() (kind SortKind, err error) {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	kind = SortKind(C.Z3_get_sort_kind(sort.context.z3val, z3sort))
	err = getError(sort.context)
	return
}

func (sort *Sort) BVSize() (size uint, err error) {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3size, err := C.Z3_get_bv_sort_size(sort.context.z3val, z3sort), getError(sort.context)
	if err != nil {
		return
	}
	size = uint(z3size)
	return
}

func (sort *Sort) ArrayDomain() (domSort *Sort, err error) {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3DomSort, err := C.Z3_get_array_sort_domain(sort.context.z3val, z3sort), getError(sort.context)
	if err != nil {
		return
	}
	domSort = newSort(sort.context, z3DomSort)
	return
}

func (sort *Sort) ArrayRange() (rangeSort *Sort, err error) {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	z3RangeSort, err := C.Z3_get_array_sort_range(sort.context.z3val, z3sort), getError(sort.context)
	if err != nil {
		return
	}
	rangeSort = newSort(sort.context, z3RangeSort)
	return
}
