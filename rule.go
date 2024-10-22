package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type RuleStatus int

const (
	Waiting RuleStatus = iota
	Completed
	Cached
	Error
)

type Rule struct {
	target       string
	dependencies []*Rule
	actions      []*Action
	time         int
	phony        bool
	Status       RuleStatus
}

func getLastModifiedTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func (r *Rule) executeActions() error {
	for _, action := range r.actions {
		if err := action.Execute(); err != nil {
			return fmt.Errorf("error executing action for rule %s: %s\n      %v", r.target, strings.Join(action.args, " "), err)
		}
	}
	return nil
}

func (r *Rule) Evaluate() error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(r.dependencies))

	for _, subRule := range r.dependencies {
		wg.Add(1)

		go func(rule *Rule) {
			defer wg.Done()
			if err := rule.Evaluate(); err != nil {
				errChan <- fmt.Errorf("error evaluating rule %s:\n    %v", rule.target, err)
			}
		}(subRule)
	}

	wg.Wait()
	close(errChan) // Close the error channel after all goroutines complete

	// Check for errors from the goroutines
	for err := range errChan {
		if err != nil {
			r.Status = Error
			return err
		}
	}

	targetModifiedTime, err := getLastModifiedTime(r.target)
	if err != nil {
		// If this returned an error, file does not exist
		if execErr := r.executeActions(); execErr != nil {
			r.Status = Error
			return execErr
		}
		r.Status = Completed
		return nil
	}

	for _, dep := range r.dependencies {
		depModTime, err := getLastModifiedTime(dep.target)
		if err != nil {
			r.Status = Error
			return fmt.Errorf("failed to get modification time of dependency %s: %v", dep.target, err)
		}

		if depModTime.After(targetModifiedTime) {
			if execErr := r.executeActions(); execErr != nil {
				r.Status = Error
				return execErr // Return the error from executeActions
			}
			r.Status = Completed
			return nil
		}
	}

	r.Status = Cached
	return nil
}
