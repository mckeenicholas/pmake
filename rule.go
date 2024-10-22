package main

import (
	"fmt"
	"sync"
)

type Rule struct {
	target       string
	dependencies []*Rule
	actions      []*Action
	time         int
	phony        bool
	completed    bool
}

func (r *Rule) Evaluate() error {
	var wg sync.WaitGroup
	var mu sync.Mutex // To handle the completed status safely

	// Evaluate each dependency in parallel
	for _, subRule := range r.dependencies {
		wg.Add(1) // Increment the WaitGroup counter

		go func(rule *Rule) {
			defer wg.Done() // Decrement the counter when done
			if err := rule.Evaluate(); err != nil {
				// Handle the error as needed (log it, return it, etc.)
				fmt.Printf("Error evaluating rule %s: %v\n", rule.target, err)
			}
			mu.Lock()
			defer mu.Unlock()
			// You can set some state here if necessary
		}(subRule)
	}

	// Wait for all dependencies to complete
	wg.Wait()

	// Execute the actions associated with the rule
	for _, action := range r.actions {
		action.Execute()
	}

	r.completed = true // Mark the rule as completed
	return nil
}
