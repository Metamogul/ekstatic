package ekstatic

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWorkflow(t *testing.T) {
	t.Parallel()

	testWorkflow := NewWorkflow()
	require.NotEmpty(t, testWorkflow)
}

func TestWorkflow_AddTransition(t *testing.T) {
	type workflow struct {
		transitions map[transitionIdentifer]Transition
	}

	testcases := []struct {
		name             string
		workflow         workflow
		transition       Transition
		initialState     any
		input            any
		destinationState any
		wantsPanicError  error
	}{
		{
			name:            "workflow not initialized",
			workflow:        workflow{},
			transition:      func(string) (string, error) { return "", nil },
			wantsPanicError: errors.New("assignment to entry in nil map"),
		},
		{
			name:            "transition is nil",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      nil,
			wantsPanicError: ErrTransitionNil,
		},
		{
			name:            "tried to add non-function",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      "foo",
			wantsPanicError: ErrTransitionIsNonFunc,
		},
		{
			name:            "transition accepts no arguments",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      func() {},
			wantsPanicError: ErrTransitionAcceptsNoArguments,
		},
		{
			name:            "transition does't have a return value",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      func(string) {},
			wantsPanicError: ErrTransitionHasNoReturnValues,
		},
		{
			name:            "transition has more than two return values",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      func(string, string) (string, string, error) { return "", "", nil },
			wantsPanicError: ErrTransitionTooManyReturnValues,
		},
		{
			name:            "second return value of transition is not an error",
			workflow:        workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition:      func(string, string) (string, string) { return "", "" },
			wantsPanicError: ErrTransitionBadErrorOutput,
		},
		{
			name: "transition with that signature already exists",
			workflow: workflow{transitions: map[transitionIdentifer]Transition{
				identifierFromTransition(func(string) (string, error) { return "", nil }): func(string) (string, error) { return "", nil },
			}},
			transition:      func(string) (string, error) { return "", nil },
			wantsPanicError: ErrTransitionAlreadyExists,
		},
		{
			name:     "transition successfully added",
			workflow: workflow{transitions: make(map[transitionIdentifer]Transition, 0)},
			transition: func(state string, input string) string {
				if state == "ekstatic" && input == "make awesome" {
					return state + " is awesome"
				} else {
					return state + " is bullshit"
				}

			},
			initialState:     "ekstatic",
			input:            "make awesome",
			destinationState: "ekstatic is awesome",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				transitions: tt.workflow.transitions,
			}

			if tt.wantsPanicError != nil {
				require.PanicsWithError(t, tt.wantsPanicError.Error(), func() { w.AddTransition(tt.transition) })
				return
			}

			w.AddTransition(tt.transition)
			instance := w.New(tt.initialState)
			err := instance.ContinueWith(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.destinationState, instance.CurrentState())
		})
	}
}

func TestWorkflow_AddTransitionSucceededAction(t *testing.T) {
	t.Parallel()

	w := NewWorkflow()
	w.AddTransition(func(state string, input string) string {
		if state == "ekstatic" && input == "make awesome" {
			return state + " is awesome"
		} else {
			return state + " is bullshit"
		}

	})

	initialState := "ekstatic"
	input := "make awesome"
	destinationState := "ekstatic is awesome"

	w.AddTransitionSucceededAction(func(newState, previousState any, transitionInput ...any) {
		require.Equal(t, initialState, previousState)
		require.Equal(t, destinationState, newState)
		require.Equal(t, input, transitionInput[0])
	})

	instance := w.New(initialState)
	err := instance.ContinueWith(input)
	require.NoError(t, err)
}

func TestWorkflow_AddTransitionFailedAction(t *testing.T) {
	t.Parallel()

	w := NewWorkflow()
	w.AddTransition(func(state string, input string) (string, error) {
		return "", errors.New("Could not make awesome")
	})

	initialState := "ekstatic"
	input := "make awesome"

	w.AddTransitionFailedAction(func(err error, previousState any, transitionInput ...any) {
		require.Equal(t, err, errors.New("Could not make awesome"))
		require.Equal(t, previousState, initialState)
		require.Equal(t, input, transitionInput[0])
	})

	instance := w.New(initialState)
	err := instance.ContinueWith(input)
	require.Error(t, err)
}

func TestWorkflow_New(t *testing.T) {
	t.Parallel()

	w := NewWorkflow()
	w.AddTransition(func(string, string) (string, error) { return "", nil })

	instance := w.New("")
	require.NotNil(t, instance)
	require.Equal(t, instance.workflow, w)
}

func TestWorkflowInstance_ContinueWith(t *testing.T) {
	t.Parallel()

	type epsilonState string

	testcases := []struct {
		name             string
		transitions      []Transition
		initialState     any
		input1           any
		input2           any
		secondInput      any
		destinationState any
		err              error
	}{
		{
			name:         "simple transition",
			initialState: "Hello, ",
			transitions: []Transition{func(state string, input string) (outputState string) {
				return state + input
			}},
			input1:           "World!",
			destinationState: "Hello, World!",
		},
		{
			name:         "simple transition with error",
			initialState: "Hello, ",
			transitions: []Transition{func(string, string) (outputState string, err error) {
				return "", errors.New("failed regular")
			}},
			input1:           "World!",
			destinationState: "Hello, ",
			err:              errors.New("failed regular"),
		},
		{
			name:         "variadic transition",
			initialState: "Hello",
			transitions: []Transition{func(state string, input string, delimeter string) (outputState string) {
				return state + delimeter + input
			}},
			input1:           "World!",
			input2:           ", ",
			destinationState: "Hello, World!",
		},
		{
			name:         "variadic transition with error",
			initialState: "Hello",
			transitions: []Transition{func(string, string, string) (outputState string, err error) {
				return "", errors.New("failed variadic")
			}},
			input1:           "World!",
			input2:           ", ",
			destinationState: "Hello",
			err:              errors.New("failed variadic"),
		},
		{
			name:         "epsilon transition",
			initialState: "Hello",
			transitions: []Transition{
				func(state string, input string) (outputState epsilonState) {
					return epsilonState(state + input)
				},
				func(state epsilonState) (outputState string) {
					return string(state) + "World!"
				},
			},
			input1:           ", ",
			destinationState: "Hello, World!",
		},
		{
			name:         "epsilon transition with error",
			initialState: "Hello",
			transitions: []Transition{
				func(state string, input string) (outputState epsilonState) {
					return epsilonState(state + input)
				},
				func(epsilonState) (outputState string, err error) {
					return "", errors.New("failed epsilon")
				},
			},
			input1:           ", ",
			destinationState: epsilonState("Hello, "),
			err:              errors.New("failed epsilon"),
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			w := NewWorkflow()
			for _, transition := range tt.transitions {
				w.AddTransition(transition)
			}

			instance := w.New(tt.initialState)

			if tt.input2 == nil {
				err := instance.ContinueWith(tt.input1)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.destinationState, instance.CurrentState())
				return
			}

			err := instance.ContinueWith(tt.input1, tt.input2)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.destinationState, instance.CurrentState())
		})
	}
}

func TestWorkflowInstance_ContinueWith_concurrenct(t *testing.T) {
	t.Parallel()

	const numberOfConcurrentCalls int = 300

	w := NewWorkflow()
	w.AddTransition(func(state string, input string) string {
		return state + input
	})
	instance := w.New("")

	wg := sync.WaitGroup{}

	callContinueWith := func(n int, input string) {
		for i := 0; i < n; i++ {
			err := instance.ContinueWith(input)
			require.NoError(t, err)
		}
		wg.Done()
	}

	wg.Add(3)
	go callContinueWith(numberOfConcurrentCalls/3, "a")
	go callContinueWith(numberOfConcurrentCalls/3, "b")
	go callContinueWith(numberOfConcurrentCalls/3, "c")
	wg.Wait()

	result, isString := instance.CurrentState().(string)
	require.True(t, isString)
	require.Len(t, result, numberOfConcurrentCalls)
	require.Equal(t, strings.Count(result, "a"), numberOfConcurrentCalls/3)
	require.Equal(t, strings.Count(result, "b"), numberOfConcurrentCalls/3)
	require.Equal(t, strings.Count(result, "c"), numberOfConcurrentCalls/3)
}
