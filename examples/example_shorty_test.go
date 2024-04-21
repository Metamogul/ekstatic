package examples

import (
	"fmt"

	"github.com/metamogul/ekstatic"
)

type initial struct {
	content string
}

type shorty int

type expanded struct {
	madeUpContent string
}

type shortenInitialInput struct{}

type expandShortInput rune

func shortenInitial(i initial, s shortenInitialInput) shorty {
	fmt.Printf("Performing shortenInitial with state: %s\n", i.content)
	fmt.Printf("Return: %d\n", len(i.content))

	return shorty(len(i.content))
}

func expandShort(s shorty, e expandShortInput) expanded {
	fmt.Printf("Performing expandShort with state: %d, input: %c\n", s, e)

	var result string

	for range s {
		result += string(e)
	}

	fmt.Printf("Return: %s\n", result)

	return expanded{result}
}

func ExampleWorkflow_shorty() {
	stateMachine := ekstatic.NewWorkflow(initial{"hello"})
	stateMachine.AddTransition(shortenInitial)
	stateMachine.AddTransition(expandShort)

	stateMachine.ContinueWith(shortenInitialInput{})
	stateMachine.ContinueWith(expandShortInput('s'))

	// Output:
	// Performing shortenInitial with state: hello
	// Return: 5
	// Performing expandShort with state: 5, input: s
	// Return: sssss
}
