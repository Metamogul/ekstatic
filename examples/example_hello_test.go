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
	triggerParse    emptyInput
	triggerTrimWith string
)

func ExampleWorkflow_hello() {
	trimHelloWorkflow := ekstatic.NewWorkflow()
	trimHelloWorkflow.AddTransition(func(s stateInput, t triggerParse) stateParsed { return stateParsed{string(s)} })
	trimHelloWorkflow.AddTransition(func(s stateParsed, t triggerTrimWith) stateTrimmed {
		result := strings.TrimSuffix(s.parsed, string(t))

		return stateTrimmed(result)
	})

	helloTrimmer := trimHelloWorkflow.New(stateInput("Hello"))

	printState(helloTrimmer)
	helloTrimmer.ContinueWith(triggerParse{})
	printState(helloTrimmer)
	helloTrimmer.ContinueWith(triggerTrimWith("llo"))
	printState(helloTrimmer)

	anotherHelloTrimmer := trimHelloWorkflow.New(stateInput("Ciao"))

	printState(anotherHelloTrimmer)
	anotherHelloTrimmer.ContinueWith(triggerParse{})
	printState(anotherHelloTrimmer)
	anotherHelloTrimmer.ContinueWith(triggerTrimWith("llo"))
	printState(anotherHelloTrimmer)

	// Output:
	// stateInput: Hello
	// stateParsed: {Hello}
	// stateTrimmed: He
	// stateInput: Ciao
	// stateParsed: {Ciao}
	// stateTrimmed: Ciao
}

func printState(sm *ekstatic.WorkflowInstance) {
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
