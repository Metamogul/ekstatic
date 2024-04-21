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

func ExampleWorkflow_persistance() {
	customerDataStore := &customerDataStore{
		records: map[int]customer{1: {
			id:        1,
			firstName: "Alex",
			lastName:  "Baker",
		}},
	}
	fmt.Println("initialized customer data store")

	updateCustomerWorkflow := newUpdateCustomerWorkflow(customerDataStore)
	fmt.Printf("created new customer update workflow\n")

	customer := customerDataStore.get(1)
	customerUpdater := updateCustomerWorkflow.New(customer)
	fmt.Printf("\nspawned new customer update workflow instance\n")

	fmt.Printf("current state: %v\n", customerUpdater.CurrentState())
	customerUpdater.ContinueWith("Chris", "Hacker")
	fmt.Printf("current state: %v\n", customerUpdater.CurrentState())
	customerUpdater.ContinueWith("Superstreet", "1b", "foo")
	fmt.Printf("current state: %v\n", customerUpdater.CurrentState())

	anotherCustomer := customerDataStore.get(1)
	anotherCustomerUpdater := updateCustomerWorkflow.New(anotherCustomer)
	fmt.Printf("\nspawned new customer update workflow instance\n")

	fmt.Printf("current state: %v\n", anotherCustomerUpdater.CurrentState())
	anotherCustomerUpdater.ContinueWith("Superstreet", "1b", "12345-6789")
	fmt.Printf("current state: %v\n", anotherCustomerUpdater.CurrentState())

	// Output:
	// initialized customer data store
	// created new customer update workflow
	//
	// spawned new customer update workflow instance
	// current state: {1 Alex Baker   }
	// writing customer with id: 1
	// current state: {1 Chris Hacker   }
	// customer update failed for customer {1 Chris Hacker   } with input [Superstreet 1b foo] (reason: postcode is not valid)
	// current state: {1 Chris Hacker   }
	//
	// spawned new customer update workflow instance
	// current state: {1 Chris Hacker   }
	// writing customer with id: 1
	// current state: {1 Chris Hacker Superstreet 1b 12345-6789}
}

func newUpdateCustomerWorkflow(customerDataStore *customerDataStore) *ekstatic.Workflow {
	updateCustomerStateWorkflow := ekstatic.NewWorkflow()

	updateCustomerStateWorkflow.AddTransition(updateName)
	updateCustomerStateWorkflow.AddTransition(updateAddress)

	updateCustomerStateWorkflow.AddTransitionSucceededAction(func(newState, previousState any, input ...any) {
		customerDataStore.put(newState.(customer))
	})
	updateCustomerStateWorkflow.AddTransitionFailedAction(func(err error, previousState any, input ...any) {
		fmt.Printf(
			"customer update failed for customer %v with input %v (reason: %s)\n",
			previousState, input, err.Error(),
		)
	})

	return updateCustomerStateWorkflow
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
