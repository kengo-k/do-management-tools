package lib

import "fmt"

type Person struct {
	Name string
}

func Hello(p Person) {
	fmt.Printf("Hello, %s!\n", p.Name)
}
