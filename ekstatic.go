package ekstatic

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Transition  any
	Termination Transition

	TransitionSucceededAction func(newState, previousState any, input ...any)
	TransitionFailedAction    func(err error, previousState any, input ...any)

	transitionIdentifer any
	stateIdentifier     reflect.Type
)

var ErrTransitionDoesNotExist = errors.New("there is no transition from the current state with the given input type")

type stateMachine interface {
	AddTransition(Transition)
	AddTermination(Termination)
	AddTransitionSucceededAction(TransitionSucceededAction)
	AddTransitionFailedAction(TransitionFailedAction)

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

func (s *StateMachine) AddTransition(t Transition) {
	if t == nil {
		panic("transition must not be nil")
	}

	transitionType := reflect.TypeOf(t)
	if transitionType.Kind() != reflect.Func {
		panic("transition must be of kind func")
	}

	if transitionType.NumIn() < 1 {
		panic("transition must accept at least accept a state argument")
	}

	switch {
	case transitionType.NumOut() < 1:
		panic("transition must return at least result state")
	case transitionType.NumOut() > 2:
		panic("transition must not have more than two return values")
	case transitionType.NumOut() == 2 && transitionType.Out(1) != reflect.TypeFor[error]():
		panic("second return value of transition must be error")
	}

	identifier := identifierFromTransition(t)

	if _, transitionExists := s.transitions[identifier]; transitionExists {
		panic("there already is a transition for that state and input type")
	}

	s.transitions[identifier] = t
}

func (s *StateMachine) AddTermination(t Termination) {
	s.AddTransition(t)
	s.terminalStates[reflect.TypeOf(t).Out(0)] = struct{}{}
}

func (s *StateMachine) AddTransitionSucceededAction(onStateUpdated TransitionSucceededAction) {
	s.onTransitionSucceeded = onStateUpdated
}

func (s *StateMachine) AddTransitionFailedAction(onTransitionFailed TransitionFailedAction) {
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

	// Recursively call submachine

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

	// Select transition

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

	// Perform failure action

	if len(transitionResult) == 2 && transitionResult[1].Interface() != nil {
		err := transitionResult[1].Interface().(error)
		if s.onTransitionFailed != nil {
			s.onTransitionFailed(err, s.currentState, input...)
		}
		return err
	}

	// Perform success action & assign state

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
