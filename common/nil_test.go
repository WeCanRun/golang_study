package common

import "testing"

func TestNil(t *testing.T) {
	var s, h *string = nil, nil
	s = h
	h = s
	// panic
	//ch := *s
	//_ = ch

	var arr []*string = nil
	t.Log(len(arr))
	for _, a := range arr {
		t.Log(a)
	}
}
