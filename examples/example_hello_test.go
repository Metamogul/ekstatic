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

func ExampleWorkflow_hello() {
	sm := ekstatic.NewWorkflow(stateInput("Hello"))

	sm.AddTransition(func(s stateInput, t triggerParse) stateParsed { return stateParsed{string(s)} })
	sm.AddTransition(func(s stateParsed, t triggerTrimWith) stateTrimmed {
		result := strings.TrimSuffix(s.parsed, string(t))

		return stateTrimmed(result)
	})

	printState(sm)
	sm.ContinueWith(triggerParse{})
	printState(sm)
	sm.ContinueWith(triggerTrimWith("llo"))
	printState(sm)

	// Output:
	// stateInput: Hello
	// stateParsed: {Hello}
	// stateTrimmed: He
}

func printState(sm *ekstatic.Workflow) {
	switch state := sm.CurrentState().(type) {
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
