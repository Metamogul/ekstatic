package examples

import (
	"fmt"
	"math/rand/v2"
	"reflect"

	"github.com/metamogul/sssm/sssm2"
)

type triggerRandom struct{}

func RunExample8() {
	sm := sssm2.NewStateMachine(state1{})

	sm.AddTransition(func(s state1, t triggerRandom) (any, error) {
		if rand.IntN(2) == 0 {
			return state1{}, nil
		} else {
			return state2{}, nil
		}
	})
	sm.AddTransition(func(s state2, t triggerRandom) (any, error) {
		if rand.IntN(2) == 0 {
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
}
