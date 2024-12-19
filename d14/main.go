package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Sample
	// GRID_ROWS = 7
	// GRID_COLS = 11
	GRID_ROWS   = 103
	GRID_COLS   = 101
	NUM_SECONDS = 100
)

type grid [][]int

var world grid
var robots []*robot

type robot struct {
	startRow    int
	startCol    int
	currentRow  int
	currentCol  int
	velocityRow int
	velocityCol int
}

type quadrant struct {
	startCol int
	startRow int
	endCol   int
	endRow   int
	count    int
	label    string
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i < GRID_ROWS; i++ {
		row := []int{}
		for j := 0; j < GRID_COLS; j++ {
			row = append(row, 0)
		}
		world = append(world, row)
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		startLocation := strings.Split(strings.TrimSpace(parts[0]), ",")
		velocity := strings.Split(strings.TrimSpace(parts[1]), ",")

		startRow := aToIIgnoreError(string(startLocation[1]))
		// 2: removes the p= component
		startCol := aToIIgnoreError(string(startLocation[0][2:]))

		robot := &robot{
			startCol:    startCol,
			startRow:    startRow,
			currentCol:  startCol,
			currentRow:  startRow,
			velocityRow: aToIIgnoreError(string(velocity[1])),
			// 2: removes the v= component
			velocityCol: aToIIgnoreError(string(velocity[0][2:])),
		}

		world[startRow][startCol]++

		robots = append(robots, robot)
	}
}

func (r *robot) reset() {
	r.currentRow = r.startRow
	r.currentCol = r.startCol
}

func (g grid) print() {
	fmt.Println("======= Grid =======")
	for _, row := range g {
		for _, v := range row {
			if v == 0 {
				fmt.Print(".")
				continue
			}
			fmt.Printf("%d", v)
		}
		fmt.Println()
	}
	fmt.Println("======= End Grid =======")
}

func (g grid) printTree() {
	for _, row := range g {
		for _, v := range row {
			if v == 0 {
				fmt.Print(" ")
				continue
			}
			fmt.Printf("%d", v)
		}
		fmt.Println()
	}
}

func (g grid) generateQuadrants() []quadrant {
	middleRow := ((len(g)) / 2)
	middleCol := ((len(g[0])) / 2)

	return []quadrant{
		{
			label:    "topLeft",
			startCol: 0,
			endCol:   middleCol - 1,
			startRow: 0,
			endRow:   middleRow - 1,
			count:    0,
		},
		{
			label:    "topRight",
			startCol: middleCol + 1,
			endCol:   len(g[0]) - 1,
			startRow: 0,
			endRow:   middleRow - 1,
			count:    0,
		},
		{
			label:    "bottomLeft",
			startCol: 0,
			endCol:   middleCol - 1,
			startRow: middleRow + 1,
			endRow:   len(g) - 1,
			count:    0,
		},
		{
			label:    "bottomRight",
			startCol: middleCol + 1,
			endCol:   len(g[0]) - 1,
			startRow: middleRow + 1,
			endRow:   len(g) - 1,
			count:    0,
		},
	}
}

// If [r.velocityRow,r.velocityCol] = [-3, 2] and len(g) = 7 and len(g[x]) = 11 (7 and 11 can be any positive value)
// and r is currently at [1, 4] moving would give
// r.currentRow = ((1 + -3) + 7) % 7 = (-2 + 7) % 7 = 5 % 7 = 5
// r. currentCol = ((4 + 2) + 11) % 11 = (6 + 11) % 11 = 17 % 11 = 6
// so robot will be at [5, 6]
func (r *robot) move(g [][]int) {
	g[r.currentRow][r.currentCol]--
	r.currentRow = ((r.currentRow + r.velocityRow) + len(g)) % len(g)
	r.currentCol = ((r.currentCol + r.velocityCol) + len(g[r.currentRow])) % len(g[r.currentRow])
	g[r.currentRow][r.currentCol]++
}

func (q quadrant) printContent(w grid) {
	fmt.Printf("======= Start Quadrant %s (s: [%d, %d] e: [%d, %d]) =======\n", q.label, q.startRow, q.startCol, q.endRow, q.endCol)
	for r := q.startRow; r <= q.endRow; r++ {
		for c := q.startCol; c <= q.endCol; c++ {
			if w[r][c] == 0 {
				fmt.Print(".")
				continue
			}
			fmt.Printf("%d", w[r][c])
		}
		fmt.Println()
	}
	fmt.Printf("======= End Quadrant %s =======\n", q.label)
}

func safetyFactor(quadrants []quadrant) int {
	result := 0
	for _, quadrant := range quadrants {
		quadrant.updateRobotCount(world)
		if result == 0 && quadrant.count > 0 {
			result = quadrant.count
			continue
		}
		result *= quadrant.count
	}

	return result
}

func (q *quadrant) updateRobotCount(w grid) {
	for r := q.startRow; r <= q.endRow; r++ {
		for c := q.startCol; c <= q.endCol; c++ {
			q.count += w[r][c]
		}
	}
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (q *quadrant) overlappingRobots(w grid) bool {
	for r := q.startRow; r <= q.endRow; r++ {
		for c := q.startCol; c <= q.endCol; c++ {
			value := w[r][c]
			if value > 0 && value != 1 {
				fmt.Println("Quadrant has overlapping robots:", w[r][c])
				return true
			}
		}
	}
	return false
}

func treeFound(quadrants []quadrant, w grid) bool {
	for _, quadrant := range quadrants {
		if quadrant.overlappingRobots(w) {
			return false
		}
	}
	return false
}

func (g grid) treeFound() bool {
	for _, row := range g {
		for _, value := range row {
			if value > 1 {
				return false
			}
		}
	}
	return true
}

func bothParts() (int, int) {
	secondsElapsed := 1
	partOne := 0
	partTwo := 0
	quadrants := world.generateQuadrants()

	for secondsElapsed <= 100 {
		for _, robot := range robots {
			robot.move(world)
		}
		secondsElapsed++
	}

	partOne = safetyFactor(quadrants)

	// Not sure if this is always the case but the assumption has been made that the christmas tree is shown when 0 robots overlap with any other robot
	// also assumes the tree is never formed in part 1, these assumptions might fail for some inputs
	for {
		for _, robot := range robots {
			robot.move(world)
		}

		if world.treeFound() {
			partTwo = secondsElapsed
			world.printTree()
			break
		}

		secondsElapsed++
	}

	return partOne, partTwo
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func main() {
	defer timer()()
	partOne, partTwo := bothParts()
	fmt.Println("Part One:", partOne)
	fmt.Println("Part Two:", partTwo)
}
