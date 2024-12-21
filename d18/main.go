package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

const (
	// sample
	// BYTES_BEFORE_RESULT = 12
	// GRID_SIZE           = 7

	// actual
	BYTES_BEFORE_RESULT = 1024
	GRID_SIZE           = 70

	EMPTY     = '.'
	CORRUPTED = '#'
)

type (
	node struct {
		y      int
		x      int
		value  rune
		parent *node
	}
	grid      map[string]*node
	direction struct {
		yOffset int
		xOffset int
	}
	input []string
)

var (
	corrupted []node
	end       = node{
		y: GRID_SIZE,
		x: GRID_SIZE,
	}
	directions = [...]direction{
		// right
		{0, 1},
		//down
		{1, 0},
		//left
		{0, -1},
		//up
		{-1, 0},
	}
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	var byteY, byteX int

	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d", &byteX, &byteY)
		corrupted = append(corrupted, node{
			y:      byteY,
			x:      byteX,
			value:  CORRUPTED,
			parent: nil,
		})
	}
}

func nodeKey(y, x int) string {
	return fmt.Sprintf("%d-%d", y, x)
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (g grid) print() {
	// for _, row := range g {
	// 	for _, n := range row {
	// 		fmt.Printf("%s", string(n.value))
	// 	}
	// 	fmt.Println()
	// }
}

func (g grid) onGrid(y, x int) bool {
	return y > -1 && x > -1 && y <= GRID_SIZE && x <= GRID_SIZE
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

func (g grid) reset() {
	for _, n := range g {
		n.parent = nil
	}
}

func (n node) equal(n2 node) bool {
	return n.x == n2.x && n.y == n2.y
}

func (n *node) children(g grid, v map[string]struct{}) []*node {
	result := []*node{}
	for _, dir := range directions {
		nextY, nextX := n.y+dir.yOffset, n.x+dir.xOffset
		key := nodeKey(nextY, nextX)
		_, seen := v[key]
		if !g.onGrid(nextY, nextX) || g[key].value == CORRUPTED || seen {
			continue
		}

		g[key].parent = n
		v[key] = struct{}{}
		result = append(result, g[key])
	}
	return result
}

func (g grid) bfs() int {
	queue := []*node{g["0-0"]}
	visited := map[string]struct{}{
		"0-0": {},
	}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.equal(end) {
			return len(curr.tracePath())
		}

		for _, child := range curr.children(g, visited) {
			queue = append(queue, child)
		}
	}

	return 0
}

func generateWorld(corrupted []node) grid {
	world := make(grid)
	for y := 0; y <= GRID_SIZE; y++ {
		for x := 0; x <= GRID_SIZE; x++ {
			n := &node{
				y:      y,
				x:      x,
				parent: nil,
				value:  EMPTY,
			}

			world[nodeKey(n.y, n.x)] = n
		}
	}

	for _, c := range corrupted {
		world[nodeKey(c.y, c.x)].value = CORRUPTED
	}

	return world
}

func bothParts() (int, int, int) {
	partOne := -1
	partTwoX := -1
	partTwoY := -1

	seed := corrupted[:BYTES_BEFORE_RESULT]
	partOne = generateWorld(seed).bfs() - 1

	// Simple brute force, can likely be a lot smarter here, input isn't large enough for me to care all that much
	for i := BYTES_BEFORE_RESULT; i <= len(corrupted); i++ {
		seed := corrupted[:i]
		if generateWorld(seed).bfs() == 0 {
			// -1 because corrupted[:i] does not include i, so it means that i-1 broke the path
			corruptedNode := corrupted[i-1]
			return partOne, corruptedNode.x, corruptedNode.y
		}
	}

	return partOne, partTwoX, partTwoY
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func main() {
	defer timer()()
	partOne, firstX, firstY := bothParts()
	fmt.Println("Part One:", partOne)
	fmt.Printf("Part Two: %d,%d\n", firstX, firstY)

}
