package common

import (
	"fmt"
	"testing"
)

type S struct {
	A A
	B int
}

type A struct {
	a string
}

func TestInitStruct(t *testing.T) {
	var s1 S
	fmt.Println(s1)

	s2 := new(S)
	fmt.Println(s2)

	s3 := S{
		A: A{},
		B: 0,
	}
	fmt.Println(s3)
}

func TestStructEquals(t *testing.T) {
	// 指针类型比较地址
	t.Log(&S{A{"mem"}, 1} == &S{A{"mem"}, 1})

	// 非指针类型递归比较每个属性
	t.Log(S{A{"mem"}, 1} == S{A{"mem"}, 1})
	t.Log(S{A{"mem"}, 1} == S{A{"cpu"}, 1})
}

func TestBuildProblem(t *testing.T) {
	fmt.Println([...]string{"1"} == [...]string{"1"})
	//fmt.Println([...]string{"1"} == [...]string{"1","2"})
	//fmt.Println([]string{"1"} == []string{"1"})
}

func TestProblem(t *testing.T) {
	m := map[string]A{"index": {"a"}}

	//m["index"].a = "A //实际上有两个返回值，不能直接赋值
	a2 := m["index"]
	a2.a = "A"

	_ = []string{"a"}
	arr := []A{{a: "a"}}
	arr[0].a = "A"
	t.Log(m, arr)
}
