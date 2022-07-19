package common

import (
	"context"
	"testing"
)

func TestDefer(t *testing.T) {
	t.Log(Defer(t))
	context.Background()
}

func Defer(t *testing.T) (i int) {
	defer t.Log("first", i)

	defer func() {
		if err := recover(); err != nil {
			t.Log("recover2")
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			t.Log("recover1")
			panic("panic2")
		}
	}()

	defer t.Log("second")
	defer t.Log("third")
	panic("panic1")

	t.Log("test")

	return i + 1
}
