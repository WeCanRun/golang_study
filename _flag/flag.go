package main

import (
	"flag"
	"fmt"
)

func testFlag() {
	var name string
	var class *string
	flag.StringVar(&name, "name", "app", "name of container")
	class = flag.String("class", "app", "class of container")
	flag.Parse()
	fmt.Println("class: ", *class, ", name: ", name)

}

func testFlagSet() {
	flag.Parse()
	var name string
	set := flag.NewFlagSet("run", flag.ExitOnError)
	set.StringVar(&name, "name", "app", "name of container")
	set.StringVar(&name, "n", "app", "name of container")
	args := flag.Args()
	if len(args) <= 0 {
		fmt.Println("err: args")
	}
	set.Parse(args[1:])
	fmt.Println("name: ", name)
}

func main() {
	testFlagSet()
}
