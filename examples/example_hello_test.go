package examples

import (
	"fmt"
	"strings"

	"github.com/metamogul/ekstatic"
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

func ExampleStateMachine_hello() {
	sm := ekstatic.NewStateMachine(stateInput("Hello"))

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

	// Output:
	// stateInput: Hello
	// stateParsed: {Hello}
	// stateTrimmed: He
}

func printState(sm *ekstatic.StateMachine) {
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
