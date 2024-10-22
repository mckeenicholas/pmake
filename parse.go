package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func isComment(line string) bool {
	for _, ch := range line {
		switch ch {
		case '#':
			return true
		case '\t', ' ':
			continue
		default:
			return false
		}
	}
	return true
}

func findOrCreateRule(rules map[string]*Rule, target string) *Rule {
	rule, exists := rules[target]
	if !exists {
		rule = &Rule{
			target:       target,
			dependencies: []*Rule{},
			actions:      []*Action{},
		}
		rules[target] = rule
	}
	return rule
}

func Parse(fp *os.File) (map[string]*Rule, *Rule) {
	rules := make(map[string]*Rule) // Store rules by their target name
	var defaultRule *Rule           // Slice to maintain the order of the rules
	var currentRule *Rule           // Track the rule we're currently processing

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case isComment(line):
			// Skip comment lines
			continue
		case strings.Contains(line, ":="):
			// Variable assignment (ignored for now, could be implemented later)
			// Add code for variable handling if needed
			continue
		case strings.Contains(line, ":") && strings.Contains(line, ";"):
			// Rule with dependencies and recipe on the same line
			parts := strings.Split(line, ";")
			ruleAndDeps := strings.TrimSpace(parts[0])
			recipe := strings.TrimSpace(parts[1])

			ruleParts := strings.Split(ruleAndDeps, ":")
			target := strings.TrimSpace(ruleParts[0])
			depNames := strings.Split(ruleParts[1], " ")

			currentRule = findOrCreateRule(rules, target)

			if defaultRule == nil {
				defaultRule = currentRule
			}

			// Add dependencies
			for _, depName := range depNames {
				depName = strings.TrimSpace(depName)
				if depName != "" {
					depRule := findOrCreateRule(rules, depName)
					currentRule.dependencies = append(currentRule.dependencies, depRule)
				}
			}

			// Add the recipe action
			currentRule.actions = append(currentRule.actions, &Action{args: []string{recipe}})
		case strings.Contains(line, ":"):
			// Rule with dependencies (no recipe on the same line)
			ruleParts := strings.Split(line, ":")
			target := strings.TrimSpace(ruleParts[0])
			depNames := strings.Split(ruleParts[1], " ")

			currentRule = findOrCreateRule(rules, target)

			if defaultRule == nil {
				defaultRule = currentRule
			}

			// Add dependencies
			for _, depName := range depNames {
				depName = strings.TrimSpace(depName)
				if depName != "" {
					depRule := findOrCreateRule(rules, depName)
					currentRule.dependencies = append(currentRule.dependencies, depRule)
				}
			}
		case strings.HasPrefix(line, "\t"):
			// Recipe (indented line)
			if currentRule != nil {
				recipe := strings.TrimSpace(line)
				currentRule.actions = append(currentRule.actions, &Action{args: []string{recipe}})
			}
		}
	}

	return rules, defaultRule
}

func PrintRules(rules map[string]*Rule) {
	for _, rule := range rules {
		fmt.Printf("Rule: %s\n", rule.target)
		fmt.Printf("Dependencies: ")
		for _, dep := range rule.dependencies {
			fmt.Printf("%s ", dep.target)
		}
		fmt.Println()
		fmt.Printf("Actions:\n")
		for _, action := range rule.actions {
			fmt.Printf("\t%s\n", strings.Join(action.args, " "))
		}
		fmt.Println()
	}
}
