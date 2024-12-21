package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type towel struct {
	design string
}

var (
	// these are sorted by length as we always want longest match first
	towelPatterns = []string{}
	designs       = []string{}
	designGraph   = map[byte][]string{}
	resultCache   = map[string]int{}
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	lineIdx := 0
	resultCache = map[string]int{}
	designGraph = map[byte][]string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if lineIdx == 0 {
			for _, pattern := range strings.Split(line, ",") {
				towelPatterns = append(towelPatterns, strings.TrimSpace(pattern))
			}
		} else {
			designs = append(designs, line)
		}
		lineIdx++
	}

	// Sort by longest so we get longest match first
	sort.Slice(towelPatterns, func(i, j int) bool {
		return len(towelPatterns[i]) > len(towelPatterns[j])
	})

	// Build a 'graph' of what starting letter can reach which pattern
	// assume map is map[r] = ['rb', 'r']
	// design is 'rc'
	// then check 'rc' starts with something in map[r] in this case in does
	for _, tp := range towelPatterns {
		firstLetter := tp[0]
		designGraph[firstLetter] = append(designGraph[firstLetter], tp)
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

// for all possible patterns when we are at a 'target[0]' iterate them (longest first as this is pre sorted)
// if target starts with a valid possible pattern from 'target[0]' then recursively check by removing the pattern
// from the string and trying again
// Can't just do a simple match on prefix since we edit the string there are edge cases that overlap and produce the wrong result
// for example is valid patterns are ['wrb', 'wr','bx'] and target design is 'wrbx' and char map[w] = ['wrb', 'wr'],  map[b] = ['bx']
// If we just checked for prefix would we get
// wrb -> prefix match, check for ('x', tps)
// b -> prefix match, FAIL as 'b' is not a valid pattern
// now if we consider all patterns from the first latter of each design/sub design we would get the right answer via
// wrbx startsWtih wrb, -> target is now 'x' -> no patterns start with 'x'
// loop back to
// wrbx startsWith wr -> target is now 'bx'
// 'bx' startsWith bx -> target is now ”
// ” so valid combo
func totalCombinations(design string, tps []string) int {
	target := strings.Clone(design)
	result := 0

	//Give this design string we know the result, could be full design or partial
	if result, ok := resultCache[design]; ok {
		return result
	}

	// Since this is a recursive func, if target is "" it means everything was valid
	if len(target) == 0 {
		result++
	} else {
		// Since we need all possible combos we no longer exit if any design is possible
		// check for every possible pattern at every possible step
		// without caching this VERY slow
		for _, possiblePattern := range designGraph[target[0]] {
			if strings.HasPrefix(target, possiblePattern) {
				combos := totalCombinations(target[len(possiblePattern):], tps)
				result += combos
			}
		}
	}

	resultCache[design] = result
	return result
}

func bothParts() (int, int) {
	possibleDesigns := 0
	possibleCombinations := 0

	for _, design := range designs {
		result := totalCombinations(design, towelPatterns)
		if result > 0 {
			possibleDesigns++
		}
		possibleCombinations += result
	}
	return possibleDesigns, possibleCombinations
}

func main() {
	defer timer()()
	partOne, partTwo := bothParts()
	fmt.Println("Part One:", partOne)
	fmt.Println("Part Two:", partTwo)

}
