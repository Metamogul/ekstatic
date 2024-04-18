package ekstatic

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Transition          any
	Termination         Transition
	transitionIdentifer any
	stateIdentifier     reflect.Type
)

var ErrNotATransition = errors.New("the parameter passed is not a transition")
var ErrTransitionExists = errors.New("there is already a transition from that state type with the given input type")
var ErrTransitionDoesNotExist = errors.New("there is no transition from the current state with the given input type")

type stateMachine interface {
	AddTransition(Transition) error
	AddTermination(Termination) error
	AddTransitionSucceededAction(func(any, any, ...any))
	AddTransitionFailedAction(func(error, any, ...any))

	Apply(...any) error
	apply(...any) error

	CurrentState() any
	IsTerminated() bool
}

type StateMachine struct {
	transitions    map[transitionIdentifer]Transition
	currentState   any
	terminalStates map[stateIdentifier]struct{}

	onTransitionSucceeded func(newState, previousState any, input ...any)
	onTransitionFailed    func(err error, previousState any, input ...any)

	mu sync.Mutex
}

func NewStateMachine(initialState any) *StateMachine {
	if initialState == nil {
		panic("initial state must not be nil")
	}

	return &StateMachine{
		transitions:    make(map[transitionIdentifer]Transition),
		currentState:   initialState,
		terminalStates: make(map[stateIdentifier]struct{}),
	}
}

func (s *StateMachine) AddTransition(t Transition) error {
	if t == nil {
		panic("transition must not be nil")
	}

	transitionType := reflect.TypeOf(t)
	if transitionType.Kind() != reflect.Func {
		return ErrNotATransition
	}

	if transitionType.NumIn() < 1 {
		return ErrNotATransition
	}

	switch {
	case transitionType.NumOut() < 1:
		return ErrNotATransition
	case transitionType.NumOut() > 2:
		return ErrNotATransition
	case transitionType.NumOut() == 2 && transitionType.Out(1) != reflect.TypeFor[error]():
		return ErrNotATransition
	}

	identifier := identifierFromTransition(t)

	if _, transitionExists := s.transitions[identifier]; transitionExists {
		return ErrTransitionExists
	}

	s.transitions[identifier] = t
	return nil
}

func (s *StateMachine) AddTermination(t Termination) error {
	err := s.AddTransition(t)
	if err != nil {
		return err
	}

	s.terminalStates[reflect.TypeOf(t).Out(0)] = struct{}{}
	return nil
}

func (s *StateMachine) AddTransitionSucceededAction(onStateUpdated func(newState, previousState any, input ...any)) {
	if onStateUpdated == nil {
		panic("action must not be nil")
	}
	s.onTransitionSucceeded = onStateUpdated
}

func (s *StateMachine) AddTransitionFailedAction(onTransitionFailed func(err error, previousState any, input ...any)) {
	if onTransitionFailed == nil {
		panic("action must not be nil")
	}
	s.onTransitionFailed = onTransitionFailed
}

// Apply will apply the input to the current state of the StateMachine,
// using the transition corresponding to the type of the input and type
// of the current state.
func (s *StateMachine) Apply(input ...any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.apply(input...)
}

func (s *StateMachine) apply(input ...any) error {
	subMachine, isStateMachine := s.currentState.(stateMachine)
	if isStateMachine && !subMachine.IsTerminated() {
		err := subMachine.apply(input...)
		switch {
		case err != nil:
			return err
		case !subMachine.IsTerminated():
			return nil
		}
	}

	identifier := identifierFromArguments(s.currentState, input...)

	if _, exists := s.transitions[identifier]; !exists {
		return ErrTransitionDoesNotExist
	}

	// Perform transition

	transition := reflect.ValueOf(s.transitions[identifier])

	transitionArgs := make([]reflect.Value, 1+len(input))
	transitionArgs[0] = reflect.ValueOf(s.currentState)
	for i, inputArg := range input {
		transitionArgs[i+1] = reflect.ValueOf(inputArg)
	}

	transitionResult := transition.Call(transitionArgs)

	if transitionResult[0].Interface() == nil {
		panic("transition returned nil as result state")
	}

	if len(transitionResult) == 2 && transitionResult[1].Interface() != nil {
		err := transitionResult[1].Interface().(error)
		if s.onTransitionFailed != nil {
			s.onTransitionFailed(err, s.currentState, input...)
		}
		return err
	}

	if s.onTransitionSucceeded != nil {
		previousState := s.currentState
		s.currentState = transitionResult[0].Interface()
		s.onTransitionSucceeded(s.currentState, previousState, input...)
	} else {
		s.currentState = transitionResult[0].Interface()
	}

	// Chain ε-transition

	identifier = identifierFromArguments(s.currentState)
	if _, exists := s.transitions[identifier]; exists {
		return s.apply()
	}

	return nil
}

func (s *StateMachine) CurrentState() any {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.currentState
}

func (s *StateMachine) IsTerminated() bool {
	_, isTerminated := s.terminalStates[reflect.TypeOf(s.currentState)]
	return isTerminated
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
