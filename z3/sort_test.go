package z3

import "testing"

func getContext() *Context {
	config := NewConfig()
	return NewContext(config)
}

func TestBVSort32(t *testing.T) {
	context := getContext()
	sort, err := context.BVSort(32)
	if err != nil {
		t.Error("Expected valid sort, got", err)
	}
	astKind, err := sort.ASTKind()
	if err != nil {
		t.Error(err)
	}
	if astKind != SortAST {
		t.Error("Expected AST kind", SortAST, "got", astKind)
	}
	sortKind, err := sort.SortKind()
	if err != nil {
		t.Error(err)
	}
	if sortKind != BVSort {
		t.Error("Expected sort kind", BVSort, "got", sortKind)
	}
}

func TestBVSort0(t *testing.T) {
	context := getContext()
	_, err := context.BVSort(0)
	if err == nil {
		t.Error("Expected error, got valid sort")
	}
}
