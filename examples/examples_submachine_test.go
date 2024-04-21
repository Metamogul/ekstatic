package examples

// The following example is based on an example provided by qmuntal at
// https://github.com/qmuntal/stateless/blob/master/example_test.go and therefore
// includes the following License as required by the original author:
//
// BSD 2-Clause License
//
// Copyright (c) 2019, Quim Muntal
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"fmt"

	"github.com/metamogul/ekstatic"
)

type (
	stateOffHook emptyState
	stateRinging emptyState

	stateConnected struct {
		*ekstatic.Workflow
	}
	stateConnectedSpeaking emptyState

	stateOnHold struct {
		*ekstatic.Workflow
	}
	stateOnHoldWaiting emptyState
	stateOnHoldMuted   emptyState

	statePhoneDestroyed string
)

type (
	triggerCallDialed             string
	triggerCallConnected          emptyTrigger
	triggerLeftMessage            emptyTrigger
	triggerPlacedOnHold           emptyTrigger
	triggerTakenOffHold           emptyTrigger
	triggerPhoneHurledAgainstWall emptyTrigger
	triggerMuteMicrophone         emptyTrigger
	triggerUnmuteMicrophone       emptyTrigger
	triggerSetVolume              int
)

func Example() {
	phoneCall := ekstatic.NewWorkflow(stateOffHook{})

	phoneCall.AddTransition(func(s stateOffHook, callee triggerCallDialed) stateRinging {
		fmt.Printf("[Phone Call] placed for : [%s]\n", callee)
		return stateRinging{}
	})

	phoneCall.AddTransition(func(stateRinging, triggerCallConnected) stateConnected {
		connectedPhoneCall := ekstatic.NewWorkflow(stateConnectedSpeaking{})

		connectedPhoneCall.AddTransition(func(s stateConnectedSpeaking, volume triggerSetVolume) stateConnectedSpeaking {
			fmt.Printf("Volume set to %d!\n", volume)
			return stateConnectedSpeaking{}
		})

		connectedPhoneCall.AddTransition(func(stateConnectedSpeaking, triggerPlacedOnHold) stateOnHold {
			phoneCallOnHold := ekstatic.NewWorkflow(stateOnHoldWaiting{})

			phoneCallOnHold.AddTransition(func(stateOnHoldWaiting, triggerMuteMicrophone) stateOnHoldMuted {
				fmt.Println("Microphone muted!")
				return stateOnHoldMuted{}
			})

			phoneCallOnHold.AddTransition(func(stateOnHoldMuted, triggerUnmuteMicrophone) stateOnHoldWaiting {
				fmt.Println("Microphone unmuted!")
				return stateOnHoldWaiting{}
			})

			phoneCallOnHold.AddTransition(func(stateOnHoldWaiting, triggerTakenOffHold) ekstatic.StateTerminated {
				return ekstatic.StateTerminated{}
			})

			phoneCallOnHold.AddTransition(func(stateOnHoldWaiting, triggerPhoneHurledAgainstWall) ekstatic.StateTerminated {
				return ekstatic.StateTerminated{}
			})

			return stateOnHold{phoneCallOnHold}
		})

		connectedPhoneCall.AddTransition(func(s stateOnHold, volume triggerTakenOffHold) stateConnectedSpeaking {
			return stateConnectedSpeaking{}
		})

		connectedPhoneCall.AddTransition(func(stateOnHold, triggerPhoneHurledAgainstWall) ekstatic.StateTerminated {
			fmt.Println("[Timer:] Call ended at 11:30am")
			return ekstatic.StateTerminated{}
		})

		fmt.Println("[Timer:] Call started at 11:00am")
		return stateConnected{connectedPhoneCall}
	})

	phoneCall.AddTransition(func(stateConnected, triggerLeftMessage) stateOffHook { return stateOffHook{} })

	phoneCall.AddTransition(func(stateConnected, triggerPhoneHurledAgainstWall) statePhoneDestroyed {
		return statePhoneDestroyed("PhoneDestroyed")
	})

	phoneCall.ContinueWith(triggerCallDialed("qmuntal"))
	phoneCall.ContinueWith(triggerCallConnected{})
	phoneCall.ContinueWith(triggerSetVolume(2))
	phoneCall.ContinueWith(triggerPlacedOnHold{})
	phoneCall.ContinueWith(triggerMuteMicrophone{})
	phoneCall.ContinueWith(triggerUnmuteMicrophone{})
	phoneCall.ContinueWith(triggerTakenOffHold{})
	phoneCall.ContinueWith(triggerSetVolume(11))
	phoneCall.ContinueWith(triggerPlacedOnHold{})
	phoneCall.ContinueWith(triggerPhoneHurledAgainstWall{})
	fmt.Printf("State is %v\n", phoneCall.CurrentState())

	// Output:
	// [Phone Call] placed for : [qmuntal]
	// [Timer:] Call started at 11:00am
	// Volume set to 2!
	// Microphone muted!
	// Microphone unmuted!
	// Volume set to 11!
	// [Timer:] Call ended at 11:30am
	// State is PhoneDestroyed
}
