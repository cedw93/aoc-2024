// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// type pageRule struct {
// 	left    int
// 	right   int
// 	rawText string
// }

// var pageRules []pageRule
// var updates [][]int
// var safeUpdates = make(map[int]struct{})

// func aToIIgnoreError(s string) int {
// 	n, _ := strconv.Atoi(s)
// 	return n
// }

// func init() {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.Contains(line, "|") {
// 			parts := strings.Split(line, "|")
// 			pageRules = append(pageRules, pageRule{aToIIgnoreError(parts[0]), aToIIgnoreError(parts[1]), line})
// 		} else {
// 			if len(line) < 1 {
// 				continue
// 			}
// 			result := make([]int, 0, len(line))
// 			for _, page := range strings.Split(line, ",") {
// 				result = append(result, aToIIgnoreError(page))
// 			}
// 			updates = append(updates, result)
// 		}
// 	}
// }

// func sliceToSet(s []int) map[int]struct{} {
// 	result := make(map[int]struct{}, len(s))
// 	for _, v := range s {
// 		if _, ok := result[v]; !ok {
// 			result[v] = struct{}{}
// 		}
// 	}
// 	return result
// }

// func indexOfValue(s []int, target int) int {
// 	for i, v := range s {
// 		if v == target {
// 			return i
// 		}
// 	}

// 	return -1
// }

// func rulesForPage(allActiveRules []pageRule, page int) []pageRule {
// 	result := make([]pageRule, 0)
// 	for _, r := range allActiveRules {
// 		if r.right == page {
// 			result = append(result, r)
// 		}
// 	}
// 	return result
// }

// func brokenRules(updates []int, pageIdx int, rulesForPage []pageRule) []pageRule {
// 	result := make([]pageRule, 0)
// 	for _, ruleForPage := range rulesForPage {
// 		indexOfDependency := indexOfValue(updates, ruleForPage.left)
// 		if indexOfDependency > pageIdx {
// 			// fmt.Printf("Rule is broken for page %d because %d is after %d. Data is (%v)\n", updates[pageIdx], updates[indexOfDependency], indexOfDependency, updates)
// 			result = append(result, ruleForPage)
// 		}
// 	}

// 	// if len(result) > 0 {
// 	// 	fmt.Printf("Broken rules for %d: %v (data: %v)", updates[pageIdx], result, updates)
// 	// }
// 	return result
// }

// func medianCorrectPage(u []int, rules []pageRule) int {
// 	for i := range u {
// 		if len(brokenRules(u, i, rulesForPage(rules, u[i]))) > 0 {
// 			return 0
// 		}
// 	}

// 	return u[len(u)/2]
// }

// func swapRuleBreakers(data []int, brokenRules []pageRule) []int {
// 	result := append([]int{}, data...)
// 	for _, r := range brokenRules {
// 		lIdx := indexOfValue(data, r.left)
// 		rIdx := indexOfValue(data, r.right)
// 		result[lIdx], result[rIdx] = result[rIdx], result[lIdx]
// 		fmt.Printf("Swapped %d with %d. Data is now %v\n", data[lIdx], data[rIdx], result)
// 	}
// 	return result
// }

// func fixPage(pageIdx int, u []int, rules []pageRule) []int {
// 	data := append([]int{}, u...)
// 	activeRulesForPage := rulesForPage(rules, data[pageIdx])
// 	pageValue := data[pageIdx]
// 	currentIdx := pageIdx
// 	currentBrokenRulesForPage := brokenRules(data, pageIdx, activeRulesForPage)
// 	for len(currentBrokenRulesForPage) > 0 {
// 		data = swapRuleBreakers(data, currentBrokenRulesForPage)
// 		currentIdx = indexOfValue(data, pageValue)
// 		currentBrokenRulesForPage = brokenRules(data, currentIdx, activeRulesForPage)
// 	}
// 	return data
// }

// func medianWithCorrection(u []int, rules []pageRule) int {
// 	data := append([]int{}, u...)
// 	for i := range u {
// 		data = fixPage(i, data, rules)
// 	}

// 	return data[len(data)/2]
// }

// func partOne() int {
// 	result := 0
// 	for i, u := range updates {
// 		activeRules := make([]pageRule, 0)
// 		uAsSet := sliceToSet(u)
// 		for _, pageRule := range pageRules {
// 			if _, leftOk := uAsSet[pageRule.left]; leftOk {
// 				if _, rightOk := uAsSet[pageRule.right]; rightOk {
// 					activeRules = append(activeRules, pageRule)
// 				}
// 			}
// 		}

// 		medianPage := medianCorrectPage(u, activeRules)

// 		if medianPage == 0 {
// 			safeUpdates[i] = struct{}{}
// 		}

// 		result += medianPage
// 	}
// 	return result
// }

// func partTwo() int {
// 	result := 0
// 	for i, u := range updates {
// 		if _, ok := safeUpdates[i]; !ok {
// 			continue
// 		}

// 		activeRules := make([]pageRule, 0)
// 		uAsSet := sliceToSet(u)
// 		for _, pageRule := range pageRules {
// 			if _, leftOk := uAsSet[pageRule.left]; leftOk {
// 				if _, rightOk := uAsSet[pageRule.right]; rightOk {
// 					activeRules = append(activeRules, pageRule)
// 				}
// 			}
// 		}

// 		medianPage := medianWithCorrection(u, activeRules)

// 		fmt.Printf("Median page for %v is %d\n", u, medianPage)

// 		result += medianPage
// 	}
// 	return result
// }

// func timer() func() {
// 	start := time.Now()
// 	return func() {
// 		fmt.Printf("took %v\n", time.Since(start))
// 	}
// }

// func main() {
// 	defer timer()()
// 	fmt.Println("Part One:", partOne())
// 	fmt.Println("Part Two:", partTwo())
// }
