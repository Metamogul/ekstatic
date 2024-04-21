package examples

import (
	"fmt"
	"math/rand/v2"
	"reflect"

	"github.com/metamogul/ekstatic"
)

type state1 emptyState
type state2 emptyState

type triggerRandom emptyInput

var random = rand.New(rand.NewPCG(42, 0))

func ExampleWorkflow_random() {
	randomizedWorkflow := ekstatic.NewWorkflow()
	randomizedWorkflow.AddTransition(func(s state1, t triggerRandom) any {
		if random.IntN(2) == 0 {
			return state1{}
		} else {
			return state2{}
		}
	})
	randomizedWorkflow.AddTransition(func(s state2, t triggerRandom) any {
		if random.IntN(2) == 0 {
			return state1{}
		} else {
			return state2{}
		}
	})

	randomizer := randomizedWorkflow.New(state1{})

	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))
	randomizer.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))
	randomizer.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))
	randomizer.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))
	randomizer.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))
	randomizer.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(randomizer.CurrentState()))

	// Output:
	// examples.state1
	// examples.state2
	// examples.state1
	// examples.state2
	// examples.state2
	// examples.state2
}
