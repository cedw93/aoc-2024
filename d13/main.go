package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

var machines []machine

type button struct {
	x float64
	y float64
}

type machine struct {
	a      button
	b      button
	prizeX float64
	prizeY float64
	winner bool
	cost   int
}

const (
	// stated in the problem as the max
	MAX_BUTTON_PRESSES = 100
	PRIZE_OFFSET       = 10000000000000
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	currentMachine := machine{winner: false, cost: -1}
	var aX, aY, bX, bY, pX, pY float64
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Button A") {
			fmt.Fscanf(strings.NewReader(line), "Button A: X+%f, Y+%f", &aX, &aY)
			continue
		}

		if strings.HasPrefix(line, "Button B") {
			fmt.Fscanf(strings.NewReader(line), "Button B: X+%f, Y+%f", &bX, &bY)
			continue
		}

		if strings.HasPrefix(line, "Prize") {
			fmt.Fscanf(strings.NewReader(line), "Prize: X=%f, Y=%f", &pX, &pY)

			currentMachine.a = button{x: aX, y: aY}
			currentMachine.b = button{x: bX, y: bY}
			currentMachine.prizeX = pX
			currentMachine.prizeY = pY
			machines = append(machines, currentMachine)
			currentMachine = machine{winner: false, cost: -1}
			continue
		}
	}
}

// Given the pair of equations that must be true to 'win'
// a*m.a.x+b*m.b.x == m.prizeX
// a*m.a.y+b*m.b.y == m.prizeY
// Solves this pair of linear equations. Yay! This will have a unique solution if one exists
// The more general form is
// | Ax * i + Bx * j = prizeX |
// | Ay * i + By * j = prizeY |
// to solve for i and j where i = number of A presses and j is number of B presses
// lets make the By and Bx equal so we can isolate the i part
// | Ax * i + Bx * j = prizeX | multiply by By
// | Ay * i + By * j = prizeY | multiple by Bx
// Should give
// | Ax * By * i + Bx * By * j = prizeX * By
// | Ay * Bx * i + By * Bx * j = prizeY * Bx
// Now the 'By * Bx * j' component is the same we can subtract the two equations from each other (Bx * By * j and By * Bx * j cancel each other out)
// Ax * By * i - Ay * Bx * i = prizeX * By - prizeY * Bx
// Group this up
// (Ax * By - Ay * Bx)i = prizeX * By - prizeY * Bx
// Now divide by left hand side which gives
//
//					(PrizeX * By) - (PrizeY * Bx)
// AbuttonPresses = ---------------------------------
//						(Ax * By) - (Ay * Bx)
//
// Now we need to find j
// Ax * i + Bx * j = prizeX, subtract Ax * i
// Bx * j = prizeX - Ax * i divide by Bx to isolate j
//
//			   		  prizeX - Ax * i
//	 ButtonPresses = -------------------
//	            			Bx
func (m machine) playOptimal(offsetPrizes bool) (bool, int) {
	if offsetPrizes {
		m.prizeX += PRIZE_OFFSET
		m.prizeY += PRIZE_OFFSET
	}
	aPressed := ((m.prizeX * m.b.y) - (m.prizeY * m.b.x)) / ((m.a.x * m.b.y) - (m.a.y * m.b.x))
	bPressed := (m.prizeX - (m.a.x * aPressed)) / m.b.x

	// Check is result mod % 1 == 0, basically is it an integer result. Some solutions will have fractional answers which mean this is not viable
	// You cannot press a button a fractional number of times
	if math.Mod(aPressed, 1) == 0 && math.Mod(bPressed, 1) == 0 {
		// fmt.Printf("Won by pressing A: %f B: %f\n", aPressed, bPressed)
		return true, int((3 * aPressed) + bPressed)
	}

	return false, 0
}

func partOne() int {
	result := 0

	for _, machine := range machines {
		isWinner, cost := machine.playOptimal(false)
		if isWinner {
			result += cost
		}
	}
	return result
}

func partTwo() int {
	result := 0

	for _, machine := range machines {
		isWinner, cost := machine.playOptimal(true)
		if isWinner {
			result += cost
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
