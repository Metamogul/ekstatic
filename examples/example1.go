package examples

import (
	"fmt"

	"github.com/metamogul/sssm/sssm"
)

type (
	StartState struct {
		*sssm.StateBase
		TestField string
	}

	EndState struct {
		*sssm.StateBase
		SomeInt int
	}

	IntermediateState struct {
		*sssm.StateBase
	}
)

func RunExample1() {
	startState := StartState{sssm.NewStateBase(), "h"}
	intermediateState := IntermediateState{sssm.NewStateBase()}
	endState := EndState{sssm.NewStateBase(), 1}

	fromStartToEndTrigger := sssm.TransitionName("fromStartToEnd")
	fromStartToIntermediateTrigger := sssm.TransitionName("2")
	fromIntermediateToEndTrigger := sssm.TransitionName("fromIntermediateToEnd")
	fromEndToStartTrigger := sssm.TransitionName("fromEndToStart")

	startState.AddTransition(fromStartToEndTrigger, &endState, func(from, to sssm.State, input ...any) error {
		fromState, _ := from.(*StartState)
		toState, _ := to.(*EndState)
		fmt.Printf("Transitioning from StartState (%s) to EndState (%d)\n", fromState.TestField, toState.SomeInt)
		return nil
	})
	startState.AddTransition(fromStartToIntermediateTrigger, &intermediateState, func(from, to sssm.State, input ...any) error {
		fmt.Println("Transitioning to IntermediateState")
		return nil
	})

	intermediateState.AddTransition(fromIntermediateToEndTrigger, &endState, nil)

	endState.AddTransition(fromEndToStartTrigger, &startState, func(from, to sssm.State, input ...any) error {
		fromState, _ := from.(*EndState)
		toState, _ := to.(*StartState)
		fmt.Printf("Transitioning from EndState (%d) to StartState (%s)\n", fromState.SomeInt, toState.TestField)
		return nil
	})

	fsm, _ := sssm.NewStateMachine(&startState, &intermediateState, &endState)
	
	_ = fsm.PerformTransition(fromStartToEndTrigger)
	_ = fsm.PerformTransition(fromEndToStartTrigger)

	err := fsm.PerformTransition(fromEndToStartTrigger)
	if err != nil {
		fmt.Println(err)
	}
	_ = fsm.PerformTransition(fromStartToIntermediateTrigger)
	_ = fsm.PerformTransition(fromIntermediateToEndTrigger)
}
