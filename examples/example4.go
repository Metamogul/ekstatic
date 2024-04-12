package examples

import (
	"fmt"

	"github.com/metamogul/sssm/sssm2"
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

func shortenInitial(i initial, s shortenInitialInput) (shorty, error) {
	fmt.Printf("Performing shortenInitial with state: %s\n", i.content)
	fmt.Printf("Return: %d\n", len(i.content))

	return shorty(len(i.content)), nil
}

func expandShort(s shorty, e expandShortInput) (expanded, error) {
	fmt.Printf("Performing expandShort with state: %d, input: %c\n", s, e)

	var result string

	for range s {
		result += string(e)
	}

	fmt.Printf("Return: %s\n", result)

	return expanded{result}, nil
}

func RunExample4() {
	stateMachine := sssm2.NewStateMachine(initial{"hello"})
	stateMachine.AddTransition(shortenInitial)
	stateMachine.AddTransition(expandShort)

	stateMachine.PerformTransition(shortenInitialInput{})
	stateMachine.PerformTransition(expandShortInput('s'))
}
