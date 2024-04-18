package examples

import (
	"fmt"

	"github.com/metamogul/ekstatic"
)

type (
	stateVoid  struct{}
	stateHello string
	stateWorld string
)

type pushTrigger struct{}

func ExampleStateMachine_epsilon_transition() {
	stateMachine := ekstatic.NewStateMachine(stateVoid{})

	stateMachine.AddTransition(func(s stateVoid, p pushTrigger) stateHello { return "Hello" })
	stateMachine.AddTransition(func(s stateHello) stateWorld { return stateWorld(string(s) + ", world!") })

	fmt.Println(stateMachine.CurrentState())
	stateMachine.Apply(pushTrigger{})
	fmt.Println(stateMachine.CurrentState())

	// Output:
	// {}
	// Hello, world!
}
