package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"slices"
	"time"
)

const (
	WALL          = '#'
	REINDEER      = 'S'
	END           = 'E'
	ROTATED_SCORE = 1000
	START_DIR     = '>'
	VALID_MOVE    = '.'
)

type (
	node struct {
		y      int
		x      int
		cost   int
		parent *node
		dir    direction
	}
	grid      [][]rune
	direction struct {
		yOffset int
		xOffset int
		value   rune
	}
	priorityQueue []*node
)

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq *priorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*node))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

var directionMap = map[rune]direction{
	'>': {0, 1, '>'},
	'v': {1, 0, 'v'},
	'<': {0, -1, '<'},
	'^': {-1, 0, '^'},
}

var maze grid

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		maze = append(maze, []rune(scanner.Text()))
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (g grid) onGrid(row, col int) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

func (n *node) allowedDirs() []direction {
	switch n.dir.value {
	case '<':
		return []direction{
			directionMap['^'],
			directionMap['v'],
			directionMap['<'],
		}
	case '>':
		return []direction{
			directionMap['>'],
			directionMap['v'],
			directionMap['^'],
		}
	case 'v':
		return []direction{
			directionMap['v'],
			directionMap['>'],
			directionMap['<'],
		}
	case '^':
		return []direction{
			directionMap['^'],
			directionMap['>'],
			directionMap['<'],
		}
	default:
		return []direction{}
	}
}

func (n *node) generatePath() []*node {
	path := []*node{}
	next := n
	for next != nil {
		path = append(path, next)
		next = next.parent
	}

	return path
}

func (n *node) tracePath() []*node {
	path := []*node{}
	next := n

	for next != nil {
		path = append(path, next)
		next = next.parent
	}
	slices.Reverse(path)
	return path

}

func (g grid) findStart() *node {
	for y, row := range g {
		for x, v := range row {
			if v == REINDEER {
				return &node{
					y:      y,
					x:      x,
					cost:   0,
					parent: nil,
					dir:    directionMap[START_DIR],
				}
			}
		}
	}
	return nil
}

func (n *node) children() []*node {
	result := []*node{}
	for _, dir := range n.allowedDirs() {
		nextY, nextX := n.y+dir.yOffset, n.x+dir.xOffset

		// Only need to check wall, not need to check bounds
		if maze[nextY][nextX] == WALL || !maze.onGrid(nextY, nextX) {
			continue
		}

		// Every move costs one but if dir changes add 1000
		costDelta := 1
		if dir.value != n.dir.value {
			costDelta += 1000
		}

		result = append(result, &node{
			y:      nextY,
			x:      nextX,
			parent: n,
			cost:   n.cost + costDelta,
			dir:    dir,
		})
	}
	return result
}

func dijkstraAllPaths(start *node) (int, int) {
	queue := make(priorityQueue, 0)
	visited := make(map[string]int)
	nodesOnAnyBestPath := make(map[[2]int]struct{})
	bestScore := math.MaxUint32

	heap.Push(&queue, start)
	heap.Init(&queue)

	for len(queue) > 0 {
		curr := heap.Pop(&queue).(*node)

		key := fmt.Sprintf("%d-%d-%d", curr.y, curr.x, curr.dir.value)

		// We've seen this node before but with a lower cost, no point checking again...
		// if cost for node + direction is lower or equal to best lest explore still
		if prevCost, ok := visited[key]; ok && curr.cost > prevCost {
			continue
		}
		visited[key] = curr.cost

		if maze[curr.y][curr.x] == END {
			if curr.cost < bestScore {
				bestScore = curr.cost
			}

			if bestScore == curr.cost {
				for _, n := range curr.tracePath() {
					// map to avoid duplicates, we don't care about direct
					// We care that they have existed on a best path at SOME point
					bestKey := [2]int{n.y, n.x}
					nodesOnAnyBestPath[bestKey] = struct{}{}
				}
			}
			continue
		}

		for _, child := range curr.children() {
			heap.Push(&queue, child)
		}
	}

	return bestScore, len(nodesOnAnyBestPath)
}

func bothParts() (int, int) {
	bestScore, numNodesOnBestPath := dijkstraAllPaths(maze.findStart())
	return bestScore, numNodesOnBestPath
}

func main() {
	defer timer()()
	partOne, partTwo := bothParts()
	fmt.Println("Part One:", partOne)
	fmt.Println("Part Two:", partTwo)

}
