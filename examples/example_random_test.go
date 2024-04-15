package examples

import (
	"fmt"
	"math/rand/v2"
	"reflect"

	"github.com/metamogul/ekstatic"
)

type triggerRandom struct{}

var random = rand.New(rand.NewPCG(42, 0))

func ExampleStateMachine_random() {
	sm := ekstatic.NewStateMachine(state1{})

	sm.AddTransition(func(s state1, t triggerRandom) (any, error) {
		if random.IntN(2) == 0 {
			return state1{}, nil
		} else {
			return state2{}, nil
		}
	})
	sm.AddTransition(func(s state2, t triggerRandom) (any, error) {
		if random.IntN(2) == 0 {
			return state1{}, nil
		} else {
			return state2{}, nil
		}
	})

	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))

	// Output:
	// examples.state1
	// examples.state2
	// examples.state1
	// examples.state2
	// examples.state2
	// examples.state2
}
