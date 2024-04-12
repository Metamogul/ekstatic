package examples

import (
	"fmt"
	"strings"

	"github.com/metamogul/sssm/sssm2"
)

type (
	stateInput  string
	stateParsed struct {
		parsed string
	}
	stateTrimmed string
)

type (
	triggerParse    struct{}
	triggerTrimWith string
)

func RunExample5() {
	sm := sssm2.NewStateMachine(stateInput("Hello"))

	sm.AddTransition(func(s stateInput, t triggerParse) (stateParsed, error) {
		return stateParsed{string(s)}, nil
	})
	sm.AddTransition(func(s stateParsed, t triggerTrimWith) (stateTrimmed, error) {
		result := strings.TrimSuffix(s.parsed, string(t))

		return stateTrimmed(result), nil
	})

	printState(sm)
	sm.PerformTransition(triggerParse{})
	printState(sm)
	sm.PerformTransition(triggerTrimWith("llo"))
	printState(sm)
}

func printState(sm *sssm2.StateMachine) {
	switch state := sm.GetCurrentState().(type) {
	case stateInput:
		fmt.Println("stateInput: " + state)
	case stateParsed:
		fmt.Printf("stateParsed: %v\n", state)
	case stateTrimmed:
		fmt.Println("stateTrimmed: " + state)
	default:
		fmt.Println("Unknown state")
	}
}
