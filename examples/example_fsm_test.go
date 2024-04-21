package examples

import (
	"errors"
	"fmt"

	"github.com/metamogul/ekstatic"
)

type (
	stateFirst  emptyState
	stateSecond emptyState
	stateThird  emptyState
	stateLast   emptyState
)

type (
	triggerFirstToSecond emptyTrigger
	triggerSecondToThird emptyTrigger
	triggerSecondToFirst emptyTrigger
	triggerThirdToLast   emptyTrigger
)

var errFailed = errors.New("failed")

func ExampleWorkflow_fsm() {
	stateMachine := ekstatic.NewWorkflow(stateFirst{})

	stateMachine.AddTransition(func(stateFirst, triggerFirstToSecond) stateSecond { return stateSecond{} })
	stateMachine.AddTransition(func(stateSecond, triggerSecondToThird) stateThird { return stateThird{} })
	stateMachine.AddTransition(func(stateSecond, triggerSecondToFirst) (stateFirst, error) { return stateFirst{}, errFailed })
	stateMachine.AddTransition(func(stateThird, triggerThirdToLast) stateLast { return stateLast{} })

	printState6(stateMachine)
	stateMachine.ContinueWith(triggerFirstToSecond{})
	printState6(stateMachine)
	err := stateMachine.ContinueWith(triggerSecondToFirst{})
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printState6(stateMachine)
	stateMachine.ContinueWith(triggerSecondToThird{})
	printState6(stateMachine)
	err = stateMachine.ContinueWith(triggerSecondToThird{})
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printState6(stateMachine)
	stateMachine.ContinueWith(triggerThirdToLast{})
	printState6(stateMachine)

	// Output:
	// stateFirst
	// stateSecond
	// error: failed
	// stateSecond
	// stateThird
	// error: there is no transition from the current state with the given input type
	// stateThird
	// stateLast
}

func printState6(sm *ekstatic.Workflow) {
	switch sm.CurrentState().(type) {
	case stateFirst:
		fmt.Println("stateFirst")
	case stateSecond:
		fmt.Println("stateSecond")
	case stateThird:
		fmt.Println("stateThird")
	case stateLast:
		fmt.Println("stateLast")
	}
}
