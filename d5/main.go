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

type node struct {
	pageNumber int
	children   []int
	inEdges    int
}

func (n *node) addChild(p int) {
	n.children = append(n.children, p)
}

func (n *node) print() {
	fmt.Printf("graph[%d].children -> %v\n", n.pageNumber, n.children)
	fmt.Printf("graph[%d].inEdges -> %d\n", n.pageNumber, n.inEdges)
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
			result := []int{}
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

func printGraph(graph map[int]*node) {
	for _, v := range graph {
		v.print()
	}
}

func createdDirectedGraph(instructions []int, rules []pageRule) map[int]*node {
	graph := make(map[int]*node, len(instructions))

	// init blank graph
	for _, pn := range instructions {
		graph[pn] = &node{
			pageNumber: pn,
			children:   []int{},
			inEdges:    0,
		}
	}

	// graph[k].children should be an int slice containing the elements that must appear after k based on the rules provided
	// graph[k].inEdges is how many elements must exist before k
	// example:
	// instructions: [61,13,29] and assume rules state 61 -> 29 -> 13
	// graph[61].children -> [13 29]
	// graph[61].inEdges -> 0
	// graph[13].children -> []
	// graph[13].inEdges -> 2
	// graph[29].children -> [13]
	// graph[29].inEdges -> 1
	// if final result should be [61,29,13]
	for _, r := range rules {
		// if can be applied to the ruleset, A|B where A and B both in instructions
		if ruleApplicable(instructions, r) {
			// rule is applicable so we need to track that B must come after A
			// We also track that B has an extra in edge, so one extra instruction needs to be before it
			graph[r.left].addChild(r.right)
			graph[r.right].inEdges++
		}
	}

	return graph
}

func fixInstructionOrder(instructions []int, rules []pageRule) []int {
	graph := createdDirectedGraph(instructions, rules)
	toProcess := []*node{}
	correctedResult := []int{}

	for _, node := range graph {
		if node.inEdges == 0 {
			toProcess = append(toProcess, node)
		}
	}

	// https://en.wikipedia.org/wiki/Topological_sorting

	for len(toProcess) > 0 {
		next := toProcess[0]
		toProcess = toProcess[1:]
		// first will have no remaining in egdes to care about
		correctedResult = append(correctedResult, next.pageNumber)
		// If graph[next] has no candidates its just added to the result on next iter
		for _, candidatePage := range next.children {
			// remove one for each edge in the inEdge[candidatePage], if its 0 it must be next element
			graph[candidatePage].inEdges--
			// if inEdges is 0 then it means candidatePage is safe to append as it does not need to come after any other page
			// if its > 0 then theres still more pages before it so it'll get processed late
			if graph[candidatePage].inEdges == 0 {
				toProcess = append(toProcess, graph[candidatePage])
			}
		}
	}

	return correctedResult
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
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

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
