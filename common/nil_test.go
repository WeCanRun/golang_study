package common

import (
	"testing"
)

func TestNil(t *testing.T) {
	var s, h *string = nil, nil
	s = h
	h = s
	// panic
	//ch := *s
	//_ = ch

	// nil
	var s1 *string
	// non-nil
	var s2 = returnV(s1)

	t.Log(s1 == nil)
	t.Log(s2 == nil)

	var arr []*string = nil
	t.Log(len(arr))
	for _, a := range arr {
		t.Log(a)
	}
}

func returnV(v interface{}) interface{} {
	return v
}
