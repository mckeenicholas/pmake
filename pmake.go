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

	if defaultRule != "" {
		rule := rules[defaultRule]
		Make(rules, rule)
	} else {
		Make(rules, parsedDefaultRule)
	}

	elapsed := time.Since(start)
	fmt.Printf("Completed in %v\n", elapsed)
}
