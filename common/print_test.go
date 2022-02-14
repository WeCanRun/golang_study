package common

import (
	"fmt"
	"testing"
)

const (
	a = iota
	b
	f
)
const (
	name = 1
	c    = iota
	d
)

func TestPrint1(t *testing.T) {
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(f)
	fmt.Println(c)
	fmt.Println(d)
}

func TestPrint2(t *testing.T) {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]
	str2[1] = "new"
	t.Log(str1)
	str2 = append(str2, "z", "x", "y")
	t.Log(str1)
}
