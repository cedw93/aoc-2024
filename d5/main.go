package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type pageRule struct {
	left    int
	right   int
	rawText string
}

var pageRules []pageRule
var updates [][]int

func aToIIgnoreError(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			pageRules = append(pageRules, pageRule{aToIIgnoreError(parts[0]), aToIIgnoreError(parts[1]), line})
		} else {
			if len(line) < 1 {
				continue
			}
			result := make([]int, 0, len(line))
			for _, page := range strings.Split(line, ",") {
				result = append(result, aToIIgnoreError(page))
			}
			updates = append(updates, result)
		}
	}
}

func indexOf(s []int, target int) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

func ruleApplicable(instructions []int, rule pageRule) bool {
	return indexOf(instructions, rule.left) > -1 && indexOf(instructions, rule.right) > -1
}

func instructionsInOrder(instructions []int, rules []pageRule) bool {
	for _, rule := range rules {
		if ruleApplicable(instructions, rule) {
			indexLeft := indexOf(instructions, rule.left)
			indexRight := indexOf(instructions, rule.right)
			// fmt.Printf("Checking that %d (i: %d) is before %d (i: %d)\n", instructions[indexLeft], indexLeft, instructions[indexRight], indexRight)
			if indexLeft < 0 || indexRight < 0 {
				continue
			}

			// if L|R if L is before R then it's out of order...
			// e.g. 47|53, 47 must have a lower index in instructions than the index of 57
			if indexLeft > indexRight {
				return false
			}
		}
	}

	return true
}

func middlePage(instructions []int) int {
	return instructions[len(instructions)/2]
}

func createdDirectedGraph(instructions []int, rules []pageRule) (map[int][]int, map[int]int) {
	graph := make(map[int][]int, len(instructions))
	inEdges := make(map[int]int)

	for _, pageNum := range instructions {
		graph[pageNum] = []int{}
		inEdges[pageNum] = 0
	}

	// graph[k] should be an int slice containing the elements that must appear after k based on the rules provided
	// inEdges[k] is how many elements must exist before k
	// example:
	// graph[61] -> [13,29]
	// inEdges[61] -> 0
	// if final result should be [61,29,13]
	for _, r := range rules {
		if ruleApplicable(instructions, r) {
			graph[r.left] = append(graph[r.left], r.right)
			inEdges[r.right] += 1
		}
	}

	return graph, inEdges
}

func partOne() int {
	result := 0
	for _, u := range updates {
		if instructionsInOrder(u, pageRules) {
			result += middlePage(u)
		}
	}
	return result
}

func fixInstructionOrder(instructions []int, rules []pageRule) []int {
	graph, inEdges := createdDirectedGraph(instructions, rules)
	toProcess := []int{}
	correctedResult := []int{}

	for pageNumber, numInEdges := range inEdges {
		if numInEdges == 0 {
			toProcess = append(toProcess, pageNumber)
		}
	}

	// https://en.wikipedia.org/wiki/Topological_sorting

	for len(toProcess) > 0 {
		next := toProcess[0]
		toProcess = toProcess[1:]
		// first will have no remaining in egdes to care about
		correctedResult = append(correctedResult, next)
		for _, candidatePage := range graph[next] {
			// remove one for each edge in the inEdge[candidatePage], if its 0 it must be next element
			inEdges[candidatePage] -= 1
			if inEdges[candidatePage] == 0 {
				toProcess = append(toProcess, candidatePage)
			}
		}
	}

	return correctedResult
}

func partTwo() int {
	result := 0
	for _, u := range updates {
		if !instructionsInOrder(u, pageRules) {
			fixed := fixInstructionOrder(u, pageRules)
			result += middlePage(fixed)
		}
	}
	return result
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
