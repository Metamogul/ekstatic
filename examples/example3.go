package examples

import (
	"fmt"
	"reflect"
)

type someInterface interface {
	someFunction()
}

type someImplementing1 struct{}

func (s someImplementing1) someFunction() {
	fmt.Println("Do something")
}

type someImplementing2 struct{}

func (s someImplementing2) someFunction() {
	fmt.Println("Do something")
}

func printType(s someInterface) {
	fmt.Printf("%v\n", reflect.TypeOf(s))
	s.someFunction()
}

func RunExample3() {
	someThing1 := someImplementing1{}
	printType(someThing1)

	someThing2 := someImplementing2{}
	printType(someThing2)
}
