package ekstatic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWorkflow(t *testing.T) {
	testWorkflow := NewWorkflow()
	require.NotEmpty(t, testWorkflow)
}

func TestWorkflow_AddTransition(t *testing.T) {
	type workflow struct {
		transitions           map[transitionIdentifer]Transition
		onTransitionSucceeded func(newState, previousState any, input ...any)
		onTransitionFailed    func(err error, previousState any, input ...any)
	}

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				transitions:           tt.workflow.transitions,
				onTransitionSucceeded: tt.workflow.onTransitionSucceeded,
				onTransitionFailed:    tt.workflow.onTransitionFailed,
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

// TODO: Special tests to add for "ContinueWith":
//	 - Concurrency
//	 - Submachines
