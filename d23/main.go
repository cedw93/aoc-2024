package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type (
	connectionPair struct {
		left  string
		right string
	}
	tripletKey struct {
		one   string
		two   string
		three string
	}
	node  map[string]struct{}
	graph struct {
		edges map[string]node
	}
)

var (
	pairs          = []connectionPair{}
	network        = graph{edges: make(map[string]node)}
	totalComputers = 0
	foundCliques   = [][]string{}
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "-")
		network.addEdge(parts[0], parts[1])
		left := parts[0]
		right := parts[1]
		pairs = append(pairs, connectionPair{
			left:  left,
			right: right,
		})
	}
	foundCliques = network.findCliques()
}

func (k tripletKey) anyStartsWith(s string) bool {
	return strings.HasPrefix(k.one, s) || strings.HasPrefix(k.two, s) || strings.HasPrefix(k.three, s)
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func sortedKey(a, b, c string) tripletKey {
	parts := []string{a, b, c}
	sort.Strings(parts)
	return tripletKey{parts[0], parts[1], parts[2]}
}

func (g graph) findCliques() [][]string {
	var cliques [][]string

	currentClique := make(node)
	candidateNodes := make(node)
	excludedNodes := make(node)
	for node := range g.edges {
		candidateNodes[node] = struct{}{}
	}

	bronKerbosch(currentClique, candidateNodes, excludedNodes, g, &cliques)

	return cliques
}

// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm
// Implements the above to find the largest 'clique' which is the largest subgraph possible
func bronKerbosch(currentClique, candidateNodes, excludedNodes node, graph graph, cliques *[][]string) {
	if len(candidateNodes) == 0 && len(excludedNodes) == 0 {
		// If we get to hear this clique has been fully explored so we
		clique := []string{}
		for node := range currentClique {
			clique = append(clique, node)
		}
		*cliques = append(*cliques, clique)
		return
	}

	for node := range candidateNodes {
		newClique := copyMap(currentClique)
		newClique[node] = struct{}{}

		newCandidateNodes := mapIntersect(candidateNodes, graph.edges[node])
		newExcludedNodes := mapIntersect(excludedNodes, graph.edges[node])
		// As per the algorithm
		// recursion with
		// newClique = currentClique ⋃ node (union)
		// newCandidateNodes = candidateNodes ⋂ (nodes children) (intersect)
		// newExcludedNodes = excludedNodes ⋂ (nodes children) (intersect)
		bronKerbosch(newClique, newCandidateNodes, newExcludedNodes, graph, cliques)

		// Remove from candidates and add it to exlcuded
		delete(candidateNodes, node)
		excludedNodes[node] = struct{}{}
	}
}

// Creates an edge between node1 (n) and another node (n2)
// Has no concept of direct it's added both ways
func (g graph) addEdge(n, n2 string) {
	if g.edges[n] == nil {
		g.edges[n] = make(node)
	}
	if g.edges[n2] == nil {
		g.edges[n2] = make(node)
	}
	g.edges[n][n2] = struct{}{}
	g.edges[n2][n] = struct{}{}
}

// Deep copy of a map for recursion use
func copyMap(m node) node {
	copy := make(node)
	for k := range m {
		copy[k] = struct{}{}
	}
	return copy
}

// returns the elements in both maps
func mapIntersect(m, m2 node) node {
	result := make(node)
	for k := range m {
		if _, exists := m2[k]; exists {
			result[k] = struct{}{}
		}
	}
	return result
}

// This is not a nice solution but since we only need to find a 'triplet' of nodes
// then this is good enough
// Checks every Node and its children
// then checks the child of its children (grandchild)
// then checks the children of the grandchildren (greatGrandChild)
// if greatGrandChild == parent then we have looped back with length of 3 like so:
// parent -> child -> grandchild -> parent (great grand child)
// Could probably be rewritten to be more generic
func partOne() int {
	result := 0
	tripletCache := make(map[tripletKey]struct{})
	for parent, children := range network.edges {
		for child := range children {
			// prevents loops since N <--> A is undirected you can get there from either node
			// so they both know about each other!
			if child == parent {
				continue
			}
			for grandChild := range network.edges[child] {
				for greatGrandChild := range network.edges[grandChild] {
					if greatGrandChild == parent {
						// Uses a hash key based on the sorted order because
						// a,b,c
						// a,c,b
						// b,c,a
						// etc are all the same triangle of nodes and we don't want to double count. Sorted key for each would be
						// a,b,c => a,b,c
						// a,c,b => a,b,c
						// b,c,a => a,b,c
						// Not super effcient doing this, there's likely a better way but it works for now,
						triplet := sortedKey(parent, child, grandChild)
						if _, seen := tripletCache[triplet]; !seen {
							tripletCache[triplet] = struct{}{}
						}
					}
				}
			}
		}
	}
	for k := range tripletCache {
		if k.anyStartsWith("t") {
			result++
		}
	}

	return result
}

// This is slightly different as it's asking for the largest connected graph, there is a known algorithm for this
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm
// We found all the cliques during init so now its just a case of processing them
func partTwo() string {
	largestNetwork := []string{}
	largestNetworkSize := 0

	for _, clique := range foundCliques {
		if l := len(clique); l > largestNetworkSize {
			largestNetwork = clique
			largestNetworkSize = l
		}
	}

	sort.Strings(largestNetwork)
	return strings.Join(largestNetwork, ",")
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())

}
