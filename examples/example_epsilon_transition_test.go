package examples

import (
	"fmt"

	"github.com/metamogul/ekstatic"
)

type (
	stateVoid  emptyState
	stateHello string
	stateWorld string
)

type pushTrigger emptyInput

func ExampleWorkflow_epsilon_transition() {
	fsmWorkflow := ekstatic.NewWorkflow()
	fsmWorkflow.AddTransition(func(s stateVoid, p pushTrigger) stateHello { return "Hello" })
	fsmWorkflow.AddTransition(func(s stateHello) stateWorld { return stateWorld(string(s) + ", world!") })

	fsm := fsmWorkflow.New(stateVoid{})

	fmt.Println(fsm.CurrentState())
	fsm.ContinueWith(pushTrigger{})
	fmt.Println(fsm.CurrentState())

	// Output:
	// {}
	// Hello, world!
}
