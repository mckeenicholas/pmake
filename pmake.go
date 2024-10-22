package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var filename string
	var defaultRule string

	flag.StringVar(&filename, "f", "Makefile", "Specify the filename to open")

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		defaultRule = args[0]
	}

	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	rules, parsedDefaultRule := Parse(file)

	var ruleToEvaluate *Rule
	if defaultRule != "" {
		ruleToEvaluate = rules[defaultRule]
	} else {
		ruleToEvaluate = parsedDefaultRule
	}

	if err := Make(rules, ruleToEvaluate); err != nil {
		elapsed := time.Since(start)
		fmt.Printf("Error:\n  %v\nTerminated in %v\n", err, elapsed)
		return
	}

	// If no error, print completion message
	elapsed := time.Since(start)
	fmt.Printf("Completed in %v\n", elapsed)
}
