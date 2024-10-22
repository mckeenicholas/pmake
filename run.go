package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var WriteMutex sync.Mutex

func PrintOutput(rule *Rule, doneChan chan bool, start time.Time) {
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			elapsed := time.Since(start).Seconds()
			fmt.Print("\033[u")
			updateTimeInDependencyTree(rule, int(elapsed*10))
		}
	}
}

func updateTimeInDependencyTree(rule *Rule, time int) {
	if rule.Status == Waiting {
		rule.time = time
	}

	timeRounded := float64(rule.time) / 10

	// Overwrite just the time value on each line
	statusSymbol := "  "
	switch rule.Status {
	case Completed:
		statusSymbol = "✅"
	case Error:
		statusSymbol = "❌"
	case Cached:
		statusSymbol = "☑️"
	}
	fmt.Printf("%5.1fs %s\n", timeRounded, statusSymbol)

	// Recursively update times for dependencies
	for _, dep := range rule.dependencies {
		updateTimeInDependencyTree(dep, time)
	}
}

func printDependencyTree(rule *Rule, level int, time int) {
	indentation := strings.Repeat("  ", level)

	rule.time = time
	timeRounded := float64(time) / 10
	fmt.Printf("%5.1fs    | %s%s\n", timeRounded, indentation, rule.target)

	for _, dep := range rule.dependencies {
		printDependencyTree(dep, level+1, time)
	}

}

func Make(rules map[string]*Rule, defaultRule *Rule) error {
	doneChan := make(chan bool)
	start := time.Now()

	fmt.Print("\033[s")
	printDependencyTree(defaultRule, 0, 0)

	go PrintOutput(defaultRule, doneChan, start)

	err := defaultRule.Evaluate()

	close(doneChan)
	fmt.Print("\033[u")
	updateTimeInDependencyTree(defaultRule, 0)

	return err
}
