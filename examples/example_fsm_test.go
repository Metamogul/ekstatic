package examples

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/metamogul/ekstatic"
)

type (
	stateFirst  emptyState
	stateSecond emptyState
	stateThird  emptyState
	stateLast   emptyState
)

type (
	triggerFirstToSecond emptyInput
	triggerSecondToThird emptyInput
	triggerSecondToFirst emptyInput
	triggerThirdToLast   emptyInput
)

var errFailed = errors.New("failed")

func ExampleWorkflow_fsm() {
	fsmWorkflow := ekstatic.NewWorkflow()
	fsmWorkflow.AddTransition(func(stateFirst, triggerFirstToSecond) stateSecond { return stateSecond{} })
	fsmWorkflow.AddTransition(func(stateSecond, triggerSecondToThird) stateThird { return stateThird{} })
	fsmWorkflow.AddTransition(func(stateSecond, triggerSecondToFirst) (stateFirst, error) { return stateFirst{}, errFailed })
	fsmWorkflow.AddTransition(func(stateThird, triggerThirdToLast) stateLast { return stateLast{} })

	fsm := fsmWorkflow.New(stateFirst{})

	printFSMState(fsm)
	fsm.ContinueWith(triggerFirstToSecond{})
	printFSMState(fsm)
	err := fsm.ContinueWith(triggerSecondToFirst{})
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printFSMState(fsm)
	fsm.ContinueWith(triggerSecondToThird{})
	printFSMState(fsm)
	err = fsm.ContinueWith(triggerSecondToThird{})
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printFSMState(fsm)
	fsm.ContinueWith(triggerThirdToLast{})
	printFSMState(fsm)

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

func printFSMState(sm *ekstatic.WorkflowInstance) {
	fmt.Println(reflect.TypeOf(sm.CurrentState()).Name())
}
