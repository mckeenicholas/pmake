package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func substVars(line string, vars map[string]string) string {
	for varName, varValue := range vars {
		line = strings.ReplaceAll(line, fmt.Sprintf("$(%s)", varName), varValue)
		line = strings.ReplaceAll(line, fmt.Sprintf("${%s}", varName), varValue)
	}
	return line
}

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
			Status:       Waiting,
		}
		rules[target] = rule
	}
	return rule
}

func Parse(fp *os.File) (map[string]*Rule, *Rule) {
	vars := make(map[string]string)
	rules := make(map[string]*Rule)
	var defaultRule *Rule
	var currentRule *Rule

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := substVars(scanner.Text(), vars)
		switch {
		case isComment(line):
			continue
		case strings.HasPrefix(line, "\t"):
			if currentRule != nil {
				recipe := strings.TrimSpace(line)
				currentRule.actions = append(currentRule.actions, &Action{args: []string{recipe}})
			}
		case strings.Contains(line, "="):
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				varValue := strings.TrimSpace(parts[1])
				vars[varName] = varValue
			}

			continue
		case strings.Contains(line, ":"):
			parts := strings.Split(line, ";")
			ruleAndDeps := strings.TrimSpace(parts[0])

			ruleParts := strings.Split(ruleAndDeps, ":")
			target := strings.TrimSpace(ruleParts[0])
			depNames := strings.Split(ruleParts[1], " ")

			currentRule = findOrCreateRule(rules, target)

			if defaultRule == nil {
				defaultRule = currentRule
			}

			for _, depName := range depNames {
				depName = strings.TrimSpace(depName)
				if depName != "" {
					depRule := findOrCreateRule(rules, depName)
					currentRule.dependencies = append(currentRule.dependencies, depRule)
				}
			}

			// If line contains ";", anything after it is an action
			if len(parts) > 1 {
				recipe := strings.TrimSpace(parts[1])
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
