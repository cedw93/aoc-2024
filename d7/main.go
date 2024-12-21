package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var calibrations []*calibration

type calibration struct {
	target int
	values []int
	valid  bool
}

func aToIIgnoreError(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		values := []int{}
		for _, value := range strings.Split(parts[1], " ") {
			if len(strings.TrimSpace(value)) > 0 {
				values = append(values, aToIIgnoreError(string(value)))
			}
		}

		calibrations = append(calibrations, &calibration{aToIIgnoreError(parts[0]), values, false})
	}
}

func validOp(target int, values []int, useConcat bool) bool {
	// fmt.Printf("Target: %d, values: %v. Checking %d\n", target, values, values[len(values)-1])
	if len(values) == 1 {
		// last element to check must match our target
		return target == values[0]
	}

	// if its a multiple we need to check if remaining values * this value are valid...
	// for example take the example 3267: 81 40 27
	// if 3267 % 27 == 0 and validOp(121, [81, 40]) is true then this is valid
	if target%values[len(values)-1] == 0 && validOp(target/values[len(values)-1], values[:len(values)-1], useConcat) {
		return true
	}

	// Now we check the sum, same example we know 27 is valid multiple so checking this
	// validOp(121, [81, 40])
	// 40 is less than 121 so lets check if 121-40 = 81
	// currentOp(81, [81])
	if target > values[len(values)-1] && validOp(target-values[len(values)-1], values[:len(values)-1], useConcat) {
		return true
	}

	// part 2 requires concat, not sure if this is most optimal way but it works!
	// convert target to a string and if the length >  last elements length (num of digits) and target ends with lastele
	// then its possible to concat this one so we need to check of the other numbers are viable
	// example
	// 1234 | 3 4 34
	// 34 can be concat'd on so we need to check if 3 4 and make 12
	if useConcat {
		targetAsString := strconv.Itoa(target)
		lastAsString := strconv.Itoa(values[len(values)-1])
		// fmt.Printf("Checking concat for %d, %v. String values: target: %s, last: %s (trimmed %s)\n", target, values, targetAsString, lastAsString, strings.TrimSuffix(targetAsString, lastAsString))
		if len(targetAsString) > len(lastAsString) && strings.HasSuffix(targetAsString, lastAsString) && validOp(aToIIgnoreError(strings.TrimSuffix(targetAsString, lastAsString)), values[:len(values)-1], useConcat) {
			return true
		}
	}

	return false
}

func (c *calibration) setIfValid(useConcat bool) {
	if validOp(c.target, c.values, useConcat) {
		c.valid = true
	}
}

func partOne() int {
	result := 0

	for _, c := range calibrations {
		c.setIfValid(false)
		if c.valid {
			result += c.target
		}
	}
	return result
}

func partTwo() int {
	result := 0

	for _, c := range calibrations {
		c.setIfValid(true)
		if c.valid {
			result += c.target
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
