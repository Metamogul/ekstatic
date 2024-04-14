package sssm2

import "testing"

type (
	testState struct{}
	input1    struct{}
	input2    struct{}
)

var testFunction = func(state testState, input1 input1, input2 input2) (any, error) {
	return nil, nil
}

func BenchmarkIdentifierGeneration(b *testing.B) {
	b.Run("BenchmarkIdentifierFromTransition", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = identifierFromTransition(testFunction)
		}
	})

	b.Run("BenchmarkIdentifierFromArguments", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = identifierFromArguments(testState{}, input1{}, input2{})
		}
	})
}
