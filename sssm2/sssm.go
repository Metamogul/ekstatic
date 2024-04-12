package sssm2

// TODOs:
// - [ ] variadic functions as transitions
// - [x] transitions can error
// - [ ] persistance hooks
// - [ ] events: AddTransitionSucceededEvent, AddTransitionFailedEvent -> Event ist irgendeine Funktion die neuen state, alten state und input übergeben bekommt
// - [ ] threadsafety
// - [ ] tests + examples
// - [ ] docs
// - [ ] benchmarks

import (
	"errors"
	"reflect"
	"sync"
)

type (
	stateType reflect.Type
	inputType reflect.Type

	stateAndInputType struct {
		stateType
		inputType
	}

	Transition any
)

var ErrNotATransition = errors.New("the parameter passed is not a transition")
var ErrTransitionExists = errors.New("there is already a transition from that state type with the given input type")
var ErrTransitionDoesNotExist = errors.New("there is no transition from the current state with the given input type")

// StateMachne

type StateMachine struct {
	transitions  map[stateAndInputType]Transition
	currentState any
	mu           sync.Mutex
}

func NewStateMachine(initialState any) *StateMachine {
	return &StateMachine{
		transitions:  make(map[stateAndInputType]Transition),
		currentState: initialState,
	}
}

func (m *StateMachine) AddTransition(t Transition) error {
	transitionType := reflect.TypeOf(t)
	if transitionType.Kind() != reflect.Func {
		return ErrNotATransition
	}

	if transitionType.NumIn() != 2 {
		return ErrNotATransition
	}

	if transitionType.NumOut() != 2 {
		return ErrNotATransition
	}

	if transitionType.Out(1) != reflect.TypeFor[error]() {
		return ErrNotATransition
	}

	transitionIdentifier := stateAndInputType{
		reflect.TypeOf(t).In(0),
		reflect.TypeOf(t).In(1),
	}

	if _, transitionExists := m.transitions[transitionIdentifier]; transitionExists {
		return ErrTransitionExists
	}

	m.transitions[transitionIdentifier] = t

	return nil
}

func (m *StateMachine) PerformTransition(input any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	transitionIdentifier := stateAndInputType{
		reflect.TypeOf(m.currentState),
		reflect.TypeOf(input),
	}

	if _, exists := m.transitions[transitionIdentifier]; !exists {
		return ErrTransitionDoesNotExist
	}

	transitionValue := reflect.ValueOf(m.transitions[transitionIdentifier])
	inputValue := reflect.ValueOf(input)
	currenStateValue := reflect.ValueOf(m.currentState)

	transitionResult := transitionValue.Call([]reflect.Value{currenStateValue, inputValue})

	err := transitionResult[1].Interface()
	if err != nil {
		return err.(error)
	}

	m.currentState = transitionResult[0].Interface()

	// Plan:
	//   1. Variadische Transition
	//   2. Damit hier Epsilon-Transition callen, aber dran denken den Mutex vorher zu unlocken, ne besser: den Mutex übergeben und für die ganze Kette den gleichen verwenden

	return nil
}

func (m *StateMachine) GetCurrentState() any {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentState
}

func (m *StateMachine) SetCurrentState(state any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentState = state
}
