package _sort

import (
	"sort"
	"testing"
	"unsafe"
)

func TestSort(t *testing.T) {
	s := []string{"a", "2", "3"}
	sort.Strings(s)
	t.Log(s)

	i := []int{1, 22, 2, 5, 6}
	sort.Ints(i)
	t.Log(i)
	a := [1]bool{true}
	sizeof := unsafe.Sizeof(a)
	t.Log(sizeof)
}
