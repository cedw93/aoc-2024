package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type stoneKey struct {
	number int
	blinks int
}

var cachedResults = make(map[stoneKey]int)
var initialStoneState []int

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		stones := strings.Split(line, " ")

		for _, stone := range stones {
			initialStoneState = append(initialStoneState, aToIIgnoreError(stone))
		}
	}
}
func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func countAfterBlinking(stoneNumber int, blinks int) int {
	if count, ok := cachedResults[stoneKey{number: stoneNumber, blinks: blinks}]; ok {
		return count
	}

	if blinks == 0 {
		// no blinks so we have exactly one copy of this stone
		return 1
	}

	// if stone is 0, its replaced with stone 1 so get the count of all previous 1s
	if stoneNumber == 0 {
		return countAfterBlinking(1, blinks-1)
	}

	stoneAsString := strconv.Itoa(stoneNumber)
	if len(stoneAsString)%2 == 0 {
		leftHalf := aToIIgnoreError(stoneAsString[len(stoneAsString)/2:])
		rightHalf := aToIIgnoreError(stoneAsString[:len(stoneAsString)/2])
		newStonesCreated := countAfterBlinking(leftHalf, blinks-1) + countAfterBlinking(rightHalf, blinks-1)
		cachedResults[stoneKey{number: stoneNumber, blinks: blinks}] = newStonesCreated
		return newStonesCreated
	}

	// StoneId Must be odd in length so only one rule left to apply

	newStonesCreated := countAfterBlinking(stoneNumber*2024, blinks-1)
	cachedResults[stoneKey{number: stoneNumber, blinks: blinks}] = newStonesCreated

	return newStonesCreated

}

func partOne() int {
	result := 0

	for _, stone := range initialStoneState {
		result += countAfterBlinking(stone, 25)
	}

	return result
}

func partTwo() int {
	result := 0

	for _, stone := range initialStoneState {
		result += countAfterBlinking(stone, 75)
	}

	return result
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
