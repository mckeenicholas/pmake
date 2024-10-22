package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Define command line flags
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

	fmt.Println(defaultRule)

	// If a default rule is provided via command line, use it
	if defaultRule != "" {
		rule := rules[defaultRule]
		Make(rules, rule)
	} else {
		// Use the parsed default rule
		Make(rules, parsedDefaultRule)
	}

	elapsed := time.Since(start)
	fmt.Printf("Completed in %v\n", elapsed)
}
