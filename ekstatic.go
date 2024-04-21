package ekstatic

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Transition any

	TransitionSucceededAction func(newState, previousState any, input ...any)
	TransitionFailedAction    func(err error, previousState any, input ...any)

	StateTerminated struct{}

	transitionIdentifer any
)

var ErrTransitionNil = errors.New("transition must not be nil")
var ErrTransitionIsNonFunc = errors.New("transition must be of kind func")
var ErrTransitionAcceptsNoArguments = errors.New("transition must accept at least a state argument")
var ErrTransitionHasNoReturnValues = errors.New("transition must return at least result state")
var ErrTransitionTooManyReturnValues = errors.New("transition must not have more than two return values")
var ErrTransitionBadErrorOutput = errors.New("second return value of transition must be error")
var ErrTransitionAlreadyExists = errors.New("there already is a transition for that state and input type")

var ErrTransitionDoesNotExist = errors.New("there is no transition from the current state with the given input type")

type (
	// Workflow enables you to model workflows from FSMs to complex nondeterministic
	// processes that depend on complex state and input data.
	// 
	// Add transitions from one state to another by calling AddTransition.
	// 
	Workflow struct {
		transitions           map[transitionIdentifer]Transition
		onTransitionSucceeded func(newState, previousState any, input ...any)
		onTransitionFailed    func(err error, previousState any, input ...any)
	}

	WorkflowInstance struct {
		*Workflow
		currentState any

		mu sync.Mutex
	}

	// workflowInstance contains all methods implemented on Workflow and is needed internally
	// to identify states that are submachines.
	workflowInstance interface {
		ContinueWith(...any) error
		continueWith(...any) error

		CurrentState() any
		IsTerminated() bool
	}
)

func NewWorkflow() *Workflow {
	return &Workflow{
		transitions: make(map[transitionIdentifer]Transition, 0),
	}
}

func (w *Workflow) AddTransition(t Transition) {
	if t == nil {
		panic(ErrTransitionNil)
	}

	transitionType := reflect.TypeOf(t)
	if transitionType.Kind() != reflect.Func {
		panic(ErrTransitionIsNonFunc)
	}

	if transitionType.NumIn() < 1 {
		panic(ErrTransitionAcceptsNoArguments)
	}

	switch {
	case transitionType.NumOut() < 1:
		panic(ErrTransitionHasNoReturnValues)
	case transitionType.NumOut() > 2:
		panic(ErrTransitionTooManyReturnValues)
	case transitionType.NumOut() == 2 && transitionType.Out(1) != reflect.TypeFor[error]():
		panic(ErrTransitionBadErrorOutput)
	}

	identifier := identifierFromTransition(t)

	if _, transitionExists := w.transitions[identifier]; transitionExists {
		panic(ErrTransitionAlreadyExists)
	}

	w.transitions[identifier] = t
}

func (w *Workflow) AddTransitionSucceededAction(onStateUpdated TransitionSucceededAction) {
	w.onTransitionSucceeded = onStateUpdated
}

func (w *Workflow) AddTransitionFailedAction(onTransitionFailed TransitionFailedAction) {
	w.onTransitionFailed = onTransitionFailed
}

func (w *Workflow) New(initialState any) *WorkflowInstance {
	if initialState == nil {
		panic("initial state must not be nil")
	}

	return &WorkflowInstance{
		Workflow:     w,
		currentState: initialState,
	}
}

// ContinueWith will apply the input to the current state of the StateMachine,
// using the transition corresponding to the type of the input and type
// of the current state.
func (w *WorkflowInstance) ContinueWith(input ...any) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.continueWith(input...)
}

func (w *WorkflowInstance) continueWith(input ...any) error {

	// Recursively call submachine

	subMachine, isStateMachine := w.currentState.(workflowInstance)
	if isStateMachine && !subMachine.IsTerminated() {
		err := subMachine.continueWith(input...)
		switch {
		case err != nil:
			return err
		case !subMachine.IsTerminated():
			return nil
		}
	}

	// Select transition

	identifier := identifierFromArguments(w.currentState, input...)

	if _, exists := w.transitions[identifier]; !exists {
		return ErrTransitionDoesNotExist
	}

	// Perform transition

	transition := reflect.ValueOf(w.transitions[identifier])

	transitionArgs := make([]reflect.Value, 1+len(input))
	transitionArgs[0] = reflect.ValueOf(w.currentState)
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
		if w.onTransitionFailed != nil {
			w.onTransitionFailed(err, w.currentState, input...)
		}
		return err
	}

	// Perform success action & assign state

	if w.onTransitionSucceeded != nil {
		previousState := w.currentState
		w.currentState = transitionResult[0].Interface()
		w.onTransitionSucceeded(w.currentState, previousState, input...)
	} else {
		w.currentState = transitionResult[0].Interface()
	}

	// Chain ε-transition

	identifier = identifierFromArguments(w.currentState)
	if _, exists := w.transitions[identifier]; exists {
		return w.continueWith()
	}

	return nil
}

func (w *WorkflowInstance) CurrentState() any {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.currentState
}

func (w *WorkflowInstance) IsTerminated() bool {
	_, isTerminated := w.currentState.(StateTerminated)
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
