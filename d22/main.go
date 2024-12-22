package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	input []int
)

const (
	PRUNE_MODULUS = 16777216
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = append(input, aToIIgnoreError(scanner.Text()))
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

func calc(secret int) int {
	// multiple by 64 -> mix -> prune
	secret = ((secret << 6) ^ secret) % PRUNE_MODULUS
	// divide by 32 -> mix -> prune
	secret = ((secret / 32) ^ secret) % PRUNE_MODULUS
	// multiply by 2048 -> mix -> prune
	secret = ((secret << 11) ^ secret) % PRUNE_MODULUS
	return secret
}

func generateDeltas(secret int) ([]int, []int) {
	subSecrets := []int{}
	deltas := []int{}
	subSecrets = append(subSecrets, secret)
	for i := 1; i <= 2000; i++ {
		secret = calc(secret)
		// % 10 will get the last digit
		lastDigit := secret % 10
		subSecrets = append(subSecrets, lastDigit)
		deltas = append(deltas, lastDigit-deltas[i-1])
	}
	return subSecrets, deltas
}

func partOne() int {
	result := 0
	for _, secret := range input {
		for range 2000 {
			secret = calc(secret)
		}
		result += secret
	}
	return result
}

func partTwo() int {
	maxPossibleBanana := -1
	secretBananaCounts := map[[4]int]int{}
	for _, secretNum := range input {
		deltas := make([]int, 2000)
		// Last digit of the current 'step' in the secret process
		// as each price is the last digit of each generated secret
		// This is the 'first' price
		currentPrice := secretNum % 10

		changeSeq := map[[4]int]int{}

		for i := range 2000 {
			secretNum = calc(secretNum)
			// This is the price of secret gen i
			newPrice := secretNum % 10
			// Delta between prev - new
			delta := newPrice - currentPrice

			deltas[i] = delta
			if i >= 3 {
				// key is a 'run' of 4 deltas to a value since we care about
				// basically means this secret results in a gain/loss of this many bananas
				key := [4]int(deltas[i-3 : i+1])
				if _, ok := changeSeq[key]; !ok {
					changeSeq[key] = newPrice
				}
			}

			currentPrice = newPrice
		}
		// For every sequence of 4 secrets/deltas
		// Sum up the deltas into a top level map for getting the max later
		// so this will give the max bananas for secret
		// this will run 2000x times per secret
		for changeSeq, bananas := range changeSeq {
			secretBananaCounts[changeSeq] += bananas
		}
	}

	// Simply find the secret with the maximum bananas (banana price)
	for _, v := range secretBananaCounts {
		if v > maxPossibleBanana {
			maxPossibleBanana = v
		}
	}

	return maxPossibleBanana
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())

}
