package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"reflect"
	"slices"
	"time"
)

type (
	node struct {
		y          int
		x          int
		value      rune
		parent     *node
		actualCost int
		heuristic  int
	}
	grid [][]node
	// every x,y tracks the
	priorityQueue []node
	direction     struct {
		yOffset int
		xOffset int
	}
)

const (
	START                     = 'S'
	END                       = 'E'
	WALL                      = '#'
	TRACK                     = '.'
	GOOD_CHEAT_DELTA          = 100
	MAX_CHEAT_LENGTH_PART_ONE = 2
	MAX_CHEAT_LENGTH_PART_TWO = 20
)

var (
	racetrack  grid
	directions = [4]direction{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	}
	bestPath     []node
	lengthOfBest int
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	y := 0
	for scanner.Scan() {
		row := []node{}
		for x, char := range scanner.Text() {
			row = append(row, node{
				y:      y,
				x:      x,
				value:  char,
				parent: nil,
			})
		}
		racetrack = append(racetrack, row)
		y++
	}

	start, end := racetrack.findStartAndEnd()
	bestPath = racetrack.aStar(start, end)
	lengthOfBest = len(bestPath) - 1
}

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].heuristic < pq[j].heuristic
}

func (pq *priorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(node))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func (g grid) print() {
	for _, row := range g {
		for _, node := range row {
			fmt.Printf("%s", string(node.value))
		}
		fmt.Println()
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func (g grid) findStartAndEnd() (node, node) {
	start := node{}
	end := node{}
	for _, row := range g {
		for _, node := range row {
			if node.value == START {
				start = node
				start.actualCost = 0
			} else if node.value == END {
				end = node
			}
			// start AND end are not node{}
			if !reflect.ValueOf(start).IsZero() && !reflect.ValueOf(end).IsZero() {
				return start, end
			}
		}
	}
	return start, end
}

func (g grid) onGrid(y, x int) bool {
	return y > -1 && x > -1 && y < len(g) && x < len(g[y])
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func manhattanDist(y, y2, x, x2 int) int {
	return abs(y-y2) + abs(x-x2)
}

func (n node) children(end node, track grid) []node {
	children := []node{}

	for _, dir := range directions {
		nextY, nextX := n.y+dir.yOffset, n.x+dir.xOffset

		if !track.onGrid(nextY, nextX) || track[nextY][nextX].value == WALL {
			continue
		}

		children = append(children, node{
			y:          nextY,
			x:          nextX,
			parent:     &n,
			actualCost: n.actualCost + 1,
			heuristic:  n.actualCost + manhattanDist(end.y, end.x, nextY, nextX),
		})
	}

	return children
}

func (n node) tracePath() []node {
	path := []node{}
	next := n

	for next.parent != nil {
		path = append(path, next)
		next = *next.parent
	}
	path = append(path, next)
	slices.Reverse(path)
	return path

}

func (g grid) aStar(start, end node) []node {
	queue := make(priorityQueue, 0)
	visited := make(map[string]int)

	heap.Push(&queue, start)
	heap.Init(&queue)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		key := fmt.Sprintf("%d-%d", curr.y, curr.x)

		if _, ok := visited[key]; ok {
			continue
		}

		visited[key] = curr.actualCost

		if curr.y == end.y && curr.x == end.x {
			return curr.tracePath()
		}

		for _, child := range curr.children(end, g) {
			heap.Push(&queue, child)
		}
	}

	return []node{}
}

// NOTE: The problem/task on AoC day 20, 2024 explicitly states there is only 1 path
// and a cheat cannot create a new path to appear so we can do simple maths
// between known places on the unique path. This will fail if that constraint is not correct for any input.
func goodCheats(cheatSize, minimumSaving int) int {
	result := 0
	// Ignore last (len(bestPath)-minimumSaving) nodes as these all have a path to end less
	// the required saving. (e.g 100). No point checking those
	for i, node := range bestPath[:len(bestPath)-minimumSaving] {
		// i+1 just to save computation, result would be the same without it
		// just means we aren't checking nodes that will have a 0 or score as they are itself or in the past
		for _, candidate := range bestPath[i+1:] {
			// Each pair of nodes is checked ONCE in this loop (except the last len(bestPath)-minimumSaving)
			// and It does the following
			// For each node on the path in order, starting at Start
			// Check it's delta (manhattan) from the candidate
			// if thats within the allowed cheat length then check how far this candidate is from the end
			// Since we know the cost for every node on the path and we know the delta (the 'cheat' length )
			// We can figure our what the new length of the path if we cheat using
			// node.cost + (lengthOfBest-candidate.cost) + delta
			// which means it takes us node.cost moves to get to current node, if i can get to candidate in delta moves and
			// i know how far candidate is from the end what's the length?
			// If this is then less than lengthOfBest-minimum saving then its a good cheat, otherwise it's not worth taking
			// Example (cheatSize = 2, minimumSaving = 10)
			// node = [1][7] cost 12
			// candidate[1][9] cost 26
			// lengthOfBest = 84
			// this would give delta = 2
			// how far is candidate[1][9] from end? we know it takes 26 to get here so cost to end is simply lengthOfBest-26 = 58
			// so 12 + 58 + 2 = 72, this means that taking this cheat has a path of 72
			// this saves lengthOfBest - 72 = 12, which is greater than minimumSaving so this is a good cheat!
			// We NEVER need to do path finding again
			delta := manhattanDist(node.y, candidate.y, node.x, candidate.x)
			if delta <= cheatSize {
				distFromEnd := lengthOfBest - candidate.actualCost
				if node.actualCost+delta+distFromEnd <= lengthOfBest-minimumSaving {
					result++
				}
			}
		}
	}
	return result
}

func main() {
	defer timer()()

	fmt.Println("Part One:", goodCheats(MAX_CHEAT_LENGTH_PART_ONE, GOOD_CHEAT_DELTA))
	fmt.Println("Part Two:", goodCheats(MAX_CHEAT_LENGTH_PART_TWO, GOOD_CHEAT_DELTA))

}
