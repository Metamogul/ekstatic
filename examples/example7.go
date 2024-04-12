package examples

import (
	"fmt"
	"reflect"
	"time"

	"github.com/metamogul/sssm/sssm2"
)

type state1 struct{}
type state2 struct{}
type state3 struct{}

type trigger1to2 struct{}
type trigger2to3 struct{}

func RunExample7() {
	sm := sssm2.NewStateMachine(state1{})

	sm.AddTransition(func(s state1, t trigger1to2) (state2, error) {
		go sm.PerformTransition(trigger2to3{})
		return state2{}, nil
	})
	sm.AddTransition(func(s state2, t trigger2to3) (state3, error) { return state3{}, nil })

	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	sm.PerformTransition(trigger1to2{})
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	//sm.PerformTransition(trigger2to3{})
	//fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
	time.Sleep(time.Second)
	fmt.Println(reflect.TypeOf(sm.GetCurrentState()))
}
