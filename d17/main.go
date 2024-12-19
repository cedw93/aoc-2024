package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_OPCODE = 7
)

var registers = make([]int64, 3, 3)
var opCodes = make(map[int64]func(a, b, c, v int64) (int64, int64, int64))
var instructions = []int64{}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register A") {
			fmt.Fscanf(strings.NewReader(line), "Register A: %d", &registers[0])
		} else if strings.HasPrefix(line, "Register B") {
			fmt.Fscanf(strings.NewReader(line), "Register B: %d", &registers[1])
		} else if strings.HasPrefix(line, "Register C") {
			fmt.Fscanf(strings.NewReader(line), "Register C: %d", &registers[2])
		} else if strings.HasPrefix(line, "Program: ") {
			var nums string
			fmt.Fscanf(strings.NewReader(line), "Program: %s", &nums)
			for _, num := range strings.Split(nums, ",") {
				instructions = append(instructions, int64(aToIIgnoreError(num)))
			}
		}
	}

	opCodes[0] = func(a, b, c, v int64) (int64, int64, int64) {
		a = a >> v
		return a, b, c
	}

	opCodes[1] = func(a, b, c, v int64) (int64, int64, int64) {
		b = b ^ v
		return a, b, c
	}

	opCodes[2] = func(a, b, c, v int64) (int64, int64, int64) {
		b = (v % 8)
		return a, b, c
	}

	opCodes[4] = func(a, b, c, v int64) (int64, int64, int64) {
		b = b ^ c
		return a, b, c
	}

	opCodes[6] = func(a, b, c, v int64) (int64, int64, int64) {
		b = a >> v
		return a, b, c
	}

	opCodes[7] = func(a, b, c, v int64) (int64, int64, int64) {
		c = a >> v
		return a, b, c
	}
}

func operandValue(a, b, c, op int64) int64 {
	switch op {
	case 0, 1, 2, 3, 7:
		return op
	case 4:
		return a
	case 5:
		return b
	case 6:
		return c
	default:
		return -1
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

func processOps(a, b, c int64) []int64 {
	var output []int64
	for instructionPtr := 0; instructionPtr < len(instructions); instructionPtr += 2 {
		// items are processed in pairs of
		// opCode decides which func should be used
		// literal operand is the literal value to pass to opCode func
		// combo value is the rule applied to the literal for that value
		// for pair [0, 1]
		// opCode = 0
		// literal = 1
		// combo = 1 (based on operandValue)
		opCode := instructions[instructionPtr]
		literalOperand := instructions[instructionPtr+1]
		comboOperand := operandValue(a, b, c, literalOperand)

		if opCode > MAX_OPCODE {
			break
		}

		// if else as if opCode == 3 and 0 (or if opCode == 5) is not 0 things will break
		if opCode == 3 {
			if a != 0 {
				instructionPtr = int(literalOperand - 2)
				continue
			}
		} else if opCode == 5 {
			output = append(output, comboOperand%8)
		} else {
			targetValue := comboOperand
			if opCode == 1 {
				targetValue = literalOperand
			}
			a, b, c = opCodes[opCode](a, b, c, targetValue)
		}
	}
	return output
}

func partOne() string {
	var result []string

	for _, num := range processOps(registers[0], registers[1], registers[2]) {
		result = append(result, strconv.FormatInt(num, 10))
	}

	return strings.Join(result[:], ",")
}

func areSameInstructions(candidate []string) bool {
	for i, v := range candidate {
		val := int64(aToIIgnoreError(v))
		if val != instructions[i] {
			return false
		}
	}
	return true
}

func partTwo() int64 {
	// A cannot be 0 due to how the operations work
	currentGuess := int64(1)
	b, c := registers[1], registers[2]
	targetLength := len(instructions)
	for {
		result := processOps(int64(currentGuess), b, c)
		if len(result) == targetLength {
			// Even though deepEqual checks length, we should still only run it when we know the lengths are the same
			// otherwise its a waste of cpu to check len(targetLength) when we know its not right
			if reflect.DeepEqual(result, instructions) {
				return int64(currentGuess)
			}
			// result and target must be equal in size, but elems ard different
			// digit N only changes every 8^n cycles so add this to our current guess
			// e.g if element 5 is wrong then this will only ever change if (guess - currentCycle % 8^5) == 0
			// So we know we need to add At least that many to our counter, so we 'jump' to the next result where digit i has changed in answer
			// Since we are in reverse we find the largest gap first, then each step will be smaller (and less wasted checks)
			// result[0] will change on every run
			for i := targetLength - 1; i >= 0; i-- {
				if instructions[i] != result[i] {
					currentGuess += int64(math.Pow(8, float64(i)))
					break
				}
			}
			continue
		}

		// Double current guess since everything operates on powers of 2 (8 really)
		// Worked for my input and several sample inputs, might not work for them all
		// If this every over estimates it NEVER gets reduced
		currentGuess <<= 1
	}
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())

}
