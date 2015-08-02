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

type Sort struct {
	AST
}

func (sort *Sort) SortKind() SortKind {
	z3sort := C.Z3_sort(unsafe.Pointer(sort.z3val))
	return SortKind(C.Z3_get_sort_kind(sort.context.z3val, z3sort))
}
