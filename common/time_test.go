package common

import (
	"testing"
	"time"
)

func TestAfterWithNil(t *testing.T) {
	var tt time.Time
	t.Logf("tt:%v", tt)
}

func TestSliceWithCap(t *testing.T) {
	ints := make([]int, 0, 2)
	t.Logf("len: %d, cap: %d", len(ints), cap(ints))

	ints = append(ints, 1)
	ints = append(ints, 2)

	t.Logf("len: %d, cap: %d", len(ints), cap(ints))
	t.Log(ints)
}
