package sssm

import (
	"errors"
	"sync"
)

type (
	TransitionName          string
	StateTransformer        func(s1, s2 State, input ...any) error
	StateActivationCallback func()

	State interface {
		AddTransition(name TransitionName, target State, transformer StateTransformer) error
		GetTransition(name TransitionName) (Transition, error)
		SetActivationCallback(callback StateActivationCallback) error
		GetActivationCallback() StateActivationCallback

		Lock()
	}

	StateBase struct {
		transitionsByName  map[TransitionName]Transition
		activationCallback StateActivationCallback
		locked             bool
		mutex              sync.Mutex
	}

	Transition struct {
		Transformer StateTransformer
		Target      State
	}

	StateMachine struct {
		currentState State
		mutex        sync.Mutex
	}

	ErrTransitionFailed struct {
		Reason error

		Name  TransitionName
		From  State
		To    State
		Input []any
	}
)

var (
	ErrUndefinedTransition = errors.New("the current state doesn't have a tansition for the given name")
	ErrNoState             = errors.New("the transitions target is empty")
	ErrTransitionExists    = errors.New("a transition for the given name has already been set")
	ErrStateLocked         = errors.New("the state is locked")
)

/* ErrTransitionFailed */

func NewErrorTransitionFailed(reason error, name TransitionName, from, to State, input ...any) ErrTransitionFailed {
	return ErrTransitionFailed{
		Reason: reason,
		Name:   name,
		From:   from,
		To:     to,
		Input:  input,
	}
}

func (e ErrTransitionFailed) Unwrap() error {
	return e.Reason
}

func (e ErrTransitionFailed) Error() string {
	return e.Reason.Error()
}

/* StateBase */

func NewStateBase() *StateBase {
	newState := &StateBase{
		transitionsByName: make(map[TransitionName]Transition),
	}

	return newState
}

func (s *StateBase) AddTransition(name TransitionName, target State, transformer StateTransformer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.locked {
		return ErrStateLocked
	}

	if target == nil {
		return ErrNoState
	}

	if _, exists := s.transitionsByName[name]; exists {
		return ErrTransitionExists
	}

	s.transitionsByName[name] = Transition{
		Target:      target,
		Transformer: transformer,
	}

	return nil
}

func (s *StateBase) GetTransition(name TransitionName) (Transition, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	transition, exists := s.transitionsByName[name]
	if !exists {
		return Transition{}, ErrUndefinedTransition
	}

	return transition, nil
}

func (s *StateBase) SetActivationCallback(callback StateActivationCallback) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.locked {
		return ErrStateLocked
	}

	if s.activationCallback != nil {
		return errors.New("State already has a callback assigned")
	}

	s.activationCallback = callback
	return nil
}

func (s *StateBase) GetActivationCallback() StateActivationCallback {
	return s.activationCallback
}

func (s *StateBase) Lock() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.locked = true
}

/* FiniteStateMachine */

func NewStateMachine(initial State, states ...State) (machine *StateMachine, err error) {
	if initial == nil {
		return nil, ErrNoState
	}

	initial.Lock()
	for _, state := range states {
		state.Lock()
	}

	machine = &StateMachine{currentState: initial}
	return machine, nil
}

func (s *StateMachine) PerformTransition(name TransitionName, input ...any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	transition, err := s.currentState.GetTransition(name)
	if err != nil {
		return err
	}

	if transition.Transformer != nil {
		err = transition.Transformer(s.currentState, transition.Target, input)
		if err != nil {
			return err
		}
	}

	s.currentState = transition.Target

	activationCallback := s.currentState.GetActivationCallback()
	if activationCallback != nil {
		activationCallback()
	}

	return nil
}

func (s *StateMachine) CurrentState() State {
	return s.currentState
}
