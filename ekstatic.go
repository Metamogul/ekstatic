package ekstatic

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Transition          any
	transitionIdentifer any
)

var ErrNotATransition = errors.New("the parameter passed is not a transition")
var ErrTransitionExists = errors.New("there is already a transition from that state type with the given input type")
var ErrTransitionDoesNotExist = errors.New("there is no transition from the current state with the given input type")

// StateMachne

type StateMachine struct {
	transitions  map[transitionIdentifer]Transition
	currentState any
	mu           sync.Mutex
}

func NewStateMachine(initialState any) *StateMachine {
	return &StateMachine{
		transitions:  make(map[transitionIdentifer]Transition),
		currentState: initialState,
	}
}

func (s *StateMachine) AddTransition(t Transition) error {
	transitionType := reflect.TypeOf(t)
	if transitionType.Kind() != reflect.Func {
		return ErrNotATransition
	}

	if transitionType.NumIn() < 1 {
		return ErrNotATransition
	}

	if transitionType.NumOut() != 2 {
		return ErrNotATransition
	}

	if transitionType.Out(1) != reflect.TypeFor[error]() {
		return ErrNotATransition
	}

	identifier := identifierFromTransition(t)

	if _, transitionExists := s.transitions[identifier]; transitionExists {
		return ErrTransitionExists
	}

	s.transitions[identifier] = t

	return nil
}

func (s *StateMachine) PerformTransition(input ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.performTransition(input...)
}

func (s *StateMachine) performTransition(input ...any) error {
	identifier := identifierFromArguments(s.currentState, input...)

	if _, exists := s.transitions[identifier]; !exists {
		return ErrTransitionDoesNotExist
	}

	transition := reflect.ValueOf(s.transitions[identifier])

	transitionArgs := make([]reflect.Value, 1+len(input))
	transitionArgs[0] = reflect.ValueOf(s.currentState)
	for i, inputArg := range input {
		transitionArgs[i+1] = reflect.ValueOf(inputArg)
	}

	transitionResult := transition.Call(transitionArgs)

	err := transitionResult[1].Interface()
	if err != nil {
		return err.(error)
	}

	s.currentState = transitionResult[0].Interface()

	// Chain ε-transition
	return s.performTransition()
}

func (s *StateMachine) GetCurrentState() any {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.currentState
}

func (s *StateMachine) SetCurrentState(state any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentState = state
}

func identifierFromTransition(t Transition) transitionIdentifer {
	transitionType := reflect.TypeOf(t)
	transitionIdentifier := ""
	for i := 0; i < transitionType.NumIn(); i++ {
		transitionIdentifier += transitionType.In(i).String()
	}

	return transitionIdentifier
}

func identifierFromArguments(state any, args ...any) transitionIdentifer {
	transitionIdentifier := reflect.TypeOf(state).String()
	for _, arg := range args {
		transitionIdentifier += reflect.TypeOf(arg).String()
	}

	return transitionIdentifier
}
