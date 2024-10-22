package main

import (
	"fmt"
	"strings"
	"time"
)

func PrintOutput(rule *Rule, doneChan chan bool, start time.Time, ruleDepth int) {
	printDependencyTree(rule, 0, 0)

	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			fmt.Printf("\033[%dA", ruleDepth)
			updateTimeInDependencyTree(rule, int(elapsed*10))
		}
	}
}

func updateTimeInDependencyTree(rule *Rule, time int) {
	if !rule.completed {
		rule.time = time
	}

	timeRounded := float64(rule.time) / 10

	// Overwrite just the time value on each line
	if !rule.completed {
		fmt.Printf("%5.1fs    |\n", timeRounded)
	} else {
		fmt.Printf("%5.1fs ✅ |\n", timeRounded)
	}

	// Recursively update times for dependencies
	for _, dep := range rule.dependencies {
		updateTimeInDependencyTree(dep, time)
	}
}

func printDependencyTree(rule *Rule, level int, time int) {
	indentation := strings.Repeat("  ", level)

	if !rule.completed {
		rule.time = time
		timeRounded := float64(time) / 10
		fmt.Printf("%5.1fs    | %s%s\n", timeRounded, indentation, rule.target)
	} else {
		timeRounded := float64(rule.time) / 10
		fmt.Printf("%5.1fs ✅ | %s%s\n", timeRounded, indentation, rule.target)
	}

	for _, dep := range rule.dependencies {
		printDependencyTree(dep, level+1, time)
	}

}

func getRuleDepth(rule *Rule) int {
	depth := 1
	for _, dep := range rule.dependencies {
		depth += getRuleDepth(dep)
	}

	return depth
}

func Make(rules map[string]*Rule, defaultRule *Rule) error {
	doneChan := make(chan bool)
	start := time.Now()
	ruleDepth := getRuleDepth(defaultRule)

	go PrintOutput(defaultRule, doneChan, start, ruleDepth)
	defaultRule.Evaluate()

	close(doneChan)

	fmt.Printf("\033[%dA", ruleDepth)
	updateTimeInDependencyTree(defaultRule, 0)

	return nil
}
