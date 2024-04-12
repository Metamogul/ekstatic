package examples

import (
	"fmt"
	"strconv"

	"github.com/metamogul/sssm/sssm"
)

type (
	State1 struct {
		sssm.StateBase
		value int
	}

	State2 struct {
		sssm.StateBase
		value string
	}
)

func RunExample2() {
	state1 := &State1{*sssm.NewStateBase(), 0}
	state2 := &State2{*sssm.NewStateBase(), ""}

	const transition1to2 = sssm.TransitionName("1to2")

	state1.value = 1
	state1.AddTransition(transition1to2, state2, func(s1, s2 sssm.State, input ...any) error {
		state1 := s1.(*State1)
		state2 := s2.(*State2)

		state2.value = strconv.Itoa(state1.value)

		return nil
	})

	state2.SetActivationCallback(func() {
		fmt.Printf("Value of state2: %s\n", state2.value)
	})

	stateMachine, _ := sssm.NewStateMachine(state1, state2)
	stateMachine.PerformTransition(transition1to2)

}