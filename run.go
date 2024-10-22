package main

import (
	"fmt"
	"strings"
	"time"
)

func PrintOutput(rule *Rule, doneChan chan bool, start time.Time) {
	// Use a ticker to update the display every 0.1 seconds
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			fmt.Print("\033[H\033[2J") // Clear the screen
			elapsed := time.Since(start).Seconds()
			printDependencyTree(rule, 0, int(elapsed*10))
		}
	}
}

func printDependencyTree(rule *Rule, level int, time int) {
	indentation := strings.Repeat("  ", level)

	if !rule.completed {
		rule.time = time
		timeRounded := float64(time) / 10
		// Right justify the time with a width of 6
		fmt.Printf("%5.1fs    | %s%s\n", timeRounded, indentation, rule.target)
	} else {
		timeRounded := float64(rule.time) / 10
		// Right justify the time with a width of 6 and mark as completed
		fmt.Printf("%5.1fs ✅ | %s%s\n", timeRounded, indentation, rule.target)
	}

	// Recursively print dependencies
	for _, dep := range rule.dependencies {
		printDependencyTree(dep, level+1, time)
	}
}

func Make(rules map[string]*Rule, defaultRule *Rule) error {

	fmt.Println(defaultRule.target)

	doneChan := make(chan bool)
	start := time.Now()
	go PrintOutput(defaultRule, doneChan, start)

	defaultRule.Evaluate()

	close(doneChan)
	fmt.Print("\033[H\033[2J") // Clear the screen
	printDependencyTree(defaultRule, 0, 0)

	return nil
}
