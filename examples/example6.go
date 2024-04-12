package examples

import (
	"errors"
	"fmt"

	"github.com/metamogul/sssm/sssm2"
)

type state string

type (
	stateFirst  state
	stateSecond state
	stateThird  state
	stateLast   state
)

type trigger string

type (
	triggerFirstToSecond trigger
	triggerSecondToThird trigger
	triggerSecondToFirst trigger
	triggerThirdToLast   trigger
)

func RunExample6() {
	stateMachine := sssm2.NewStateMachine(stateFirst("initial"))

	stateMachine.AddTransition(func(s stateFirst, t triggerFirstToSecond) (stateSecond, error) { return "", nil })
	stateMachine.AddTransition(func(s stateSecond, t triggerSecondToThird) (stateThird, error) { return "", nil })
	stateMachine.AddTransition(func(s stateSecond, t triggerSecondToFirst) (stateFirst, error) { return "", errors.New("failed") })
	stateMachine.AddTransition(func(s stateSecond, t triggerSecondToThird) (stateThird, error) { return "", nil })
	stateMachine.AddTransition(func(s stateThird, t triggerThirdToLast) (stateLast, error) { return "", nil })

	printState6(stateMachine)
	stateMachine.PerformTransition(triggerFirstToSecond(""))
	printState6(stateMachine)
	err := stateMachine.PerformTransition(triggerSecondToFirst(""))
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printState6(stateMachine)
	stateMachine.PerformTransition(triggerSecondToThird(""))
	printState6(stateMachine)
	err = stateMachine.PerformTransition(triggerSecondToThird(""))
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	printState6(stateMachine)
	stateMachine.PerformTransition(triggerThirdToLast(""))
	printState6(stateMachine)
}

func printState6(sm *sssm2.StateMachine) {
	switch sm.GetCurrentState().(type) {
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
