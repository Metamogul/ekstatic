package examples

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/metamogul/ekstatic"
)

type customer struct {
	id              int
	firstName       string
	lastName        string
	streetAndNumber string
	city            string
	postcode        string
}

type customerDataStore struct {
	records map[int]customer
}

func (d *customerDataStore) get(id int) customer {
	return d.records[id]
}

func (d *customerDataStore) put(c customer) {
	fmt.Printf("writing customer with id: %d\n", c.id)
	d.records[c.id] = c
}

var testCustomerDataStore = customerDataStore{
	records: map[int]customer{1: {
		id:        1,
		firstName: "Alex",
		lastName:  "Baker",
	}},
}

func ExampleStateMachine_persistance() {
	customer := testCustomerDataStore.get(1)
	stateMachine := newUpdateCustomerStateMachine(customer)

	fmt.Printf("current state: %v\n", stateMachine.CurrentState())
	stateMachine.PerformTransition("Chris", "Hacker")
	fmt.Printf("current state: %v\n", stateMachine.CurrentState())
	stateMachine.PerformTransition("Superstreet", "1b", "foo")
	fmt.Printf("current state: %v\n", stateMachine.CurrentState())

	customer = testCustomerDataStore.get(1)
	stateMachine = newUpdateCustomerStateMachine(customer)

	fmt.Printf("current state: %v\n", stateMachine.CurrentState())
	stateMachine.PerformTransition("Superstreet", "1b", "12345-6789")
	fmt.Printf("current state: %v\n", stateMachine.CurrentState())

	// Output:
	//
	// created new state machine
	// current state: {1 Alex Baker   }
	// writing customer with id: 1
	// current state: {1 Chris Hacker   }
	// customer update failed for customer {1 Chris Hacker   } with input [Superstreet 1b foo] (reason: postcode is not valid)
	// current state: {1 Chris Hacker   }
	//
	// created new state machine
	// current state: {1 Chris Hacker   }
	// writing customer with id: 1
	// current state: {1 Chris Hacker Superstreet 1b 12345-6789}
}

func newUpdateCustomerStateMachine(c customer) *ekstatic.StateMachine {
	updateCustomerStateMachine := ekstatic.NewStateMachine(c)

	updateCustomerStateMachine.AddTransition(updateName)
	updateCustomerStateMachine.AddTransition(updateAddress)

	updateCustomerStateMachine.AddTransitionSucceededAction(func(newState, previousState any, input ...any) {
		testCustomerDataStore.put(newState.(customer))
	})
	updateCustomerStateMachine.AddTransitionFailedAction(func(err error, previousState any, input ...any) {
		fmt.Printf(
			"customer update failed for customer %v with input %v (reason: %s)\n",
			previousState, input, err.Error(),
		)
	})

	fmt.Println("")
	fmt.Println("created new state machine")

	return updateCustomerStateMachine
}

func updateName(c customer, firstName, lastName string) customer {
	c.firstName = firstName
	c.lastName = lastName

	return c
}

var errInvalidPostcode = errors.New("postcode is not valid")

func updateAddress(c customer, streetAndNumber, city, postcode string) (customer, error) {
	isValidPostCode, _ := regexp.MatchString(`\d{5}-\d{4}`, postcode)
	if !isValidPostCode {
		return customer{}, errInvalidPostcode
	}

	c.streetAndNumber = streetAndNumber
	c.city = city
	c.postcode = postcode

	return c, nil
}
