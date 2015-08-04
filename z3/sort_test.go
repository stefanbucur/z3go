package z3

import "testing"

func getContext() *Context {
	config := NewConfig()
	return NewContext(config)
}

func TestBVSort32(t *testing.T) {
	ctx := getContext()
	sort := ctx.BVSort(32)
	if sort == nil {
		t.Error("Expected valid sort, got nil")
	}
	if astKind := sort.ASTKind(); astKind != SortAST {
		t.Error("Expected AST kind", SortAST, "got", astKind)
	}
	if sortKind := sort.SortKind(); sortKind != BVSort {
		t.Error("Expected sort kind", BVSort, "got", sortKind)
	}
}

func TestBVSort0(t *testing.T) {
	ctx := getContext()

	if sort := ctx.BVSort(0); sort != nil {
		t.Error("Expected error, got valid sort")
	}
}
