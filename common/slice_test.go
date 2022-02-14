package common

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	a := new([]int)
	b := []int{1, 2, 3}
	c := a
	fmt.Println("before")
	fmt.Println(a, &a)
	fmt.Println(c, &c)
	*a = append(b, 4)
	fmt.Println("after")
	fmt.Println(a, &a)
	fmt.Println(c, &c)
}
