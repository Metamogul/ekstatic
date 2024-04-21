package examples

import (
	"fmt"
	"math/rand/v2"
	"reflect"

	"github.com/metamogul/ekstatic"
)

type triggerRandom struct{}

var random = rand.New(rand.NewPCG(42, 0))

func ExampleWorkflow_random() {
	sm := ekstatic.NewWorkflow(state1{})

	sm.AddTransition(func(s state1, t triggerRandom) any {
		if random.IntN(2) == 0 {
			return state1{}
		} else {
			return state2{}
		}
	})
	sm.AddTransition(func(s state2, t triggerRandom) any {
		if random.IntN(2) == 0 {
			return state1{}
		} else {
			return state2{}
		}
	})

	fmt.Println(reflect.TypeOf(sm.CurrentState()))
	sm.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.CurrentState()))
	sm.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.CurrentState()))
	sm.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.CurrentState()))
	sm.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.CurrentState()))
	sm.ContinueWith(triggerRandom{})
	fmt.Println(reflect.TypeOf(sm.CurrentState()))

	// Output:
	// examples.state1
	// examples.state2
	// examples.state1
	// examples.state2
	// examples.state2
	// examples.state2
}
