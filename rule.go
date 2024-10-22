package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Rule struct {
	target       string
	dependencies []*Rule
	actions      []*Action
	time         int
	phony        bool
	completed    bool
}

func getLastModifiedTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func (r *Rule) executeActions() {
	for _, action := range r.actions {
		action.Execute()
	}
}

func (r *Rule) Evaluate() error {
	var wg sync.WaitGroup

	for _, subRule := range r.dependencies {
		wg.Add(1)

		go func(rule *Rule) {
			defer wg.Done()
			if err := rule.Evaluate(); err != nil {
				fmt.Printf("Error evaluating rule %s: %v\n", rule.target, err)
			}
		}(subRule)
	}

	wg.Wait()

	targetModifiedTime, err := getLastModifiedTime(r.target)

	if err != nil {
		// If this returned an error, file does not exist
		r.executeActions()
		r.completed = true
		return nil
	}

	for _, dep := range r.dependencies {
		depModTime, err := getLastModifiedTime(dep.target)
		if err != nil {
			return fmt.Errorf("failed to get modification time of dependency %s: %v", dep.target, err)
		}

		if depModTime.After(targetModifiedTime) {
			r.executeActions()

			r.completed = true
			return nil
		}
	}

	r.completed = true
	return nil
}
