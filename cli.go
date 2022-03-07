package gocli

import (
	"fmt"
	"os"
)

func NewCli(entrypoint *Command) Cli {
	return Cli{
		Entrypoint:  entrypoint,
		childrenMap: make(map[*Command][]*Command),
	}
}

type Cli struct {
	Entrypoint *Command

	// Maps commands to their children
	childrenMap map[*Command][]*Command
}

func (cli *Cli) Exec() {
	args := os.Args[1:]

	// Check that the CLI tree structure is valid
	g := cli2Graph(cli)
	if !g.isTree() {
		panic(fmt.Errorf("Cli command tree has invalid structure."))
	}
	// Run the root
	cli.Entrypoint.RunUtil(args, cli.childrenMap, []string{})
}

func (cli *Cli) AddChild(parent *Command, child *Command) error {
	if cli.HasChild(parent, child) {
		return fmt.Errorf("Duplicate child in CLI tree. Parent: \"%s\", child: \"%s\"", parent.Name, child.Name)
	}
	children := cli.childrenMap[parent]
	cli.childrenMap[parent] = append(children, child)

	return nil
}

func (cli *Cli) HasChild(parent *Command, child *Command) bool {
	children := cli.childrenMap[parent]

	for _, c := range children {
		if c.Name == child.Name {
			return true
		}
	}

	return false
}

// convert childrenMap to Graph to check for cyclic
func cli2Graph(cli *Cli) Graph {
	edges := [][]int{}
	keys := make(map[*Command]int)
	keys[cli.Entrypoint] = 0

	k := 1
	for parent := range cli.childrenMap {
		if _, ok := keys[parent]; !ok {
			keys[parent] = k
			k++
		}
		for _, child := range cli.childrenMap[parent] {
			if _, ok := keys[child]; !ok {
				keys[child] = k
				k++
			}
			edges = append(edges, []int{keys[parent], keys[child]})
		}
	}

	g := Graph{
		Adj: make([][]bool, k, k),
	}

	for i := 0; i < k; i++ {
		g.Adj[i] = make([]bool, k)
	}

	for _, arr := range edges {
		g.Adj[arr[0]][arr[1]] = true
		// g.Adj[arr[1]][arr[0]] = true
	}

	return g
}

// undirected graph
type Graph struct {
	// adjacency matrix
	Adj [][]bool
}

func (g *Graph) isCyclic(v int, visited map[int]bool, parent int) bool {
	if len(g.Adj) == 0 {
		return false
	}

	if len(g.Adj) == 1 {
		visited[0] = true
		if len(g.Adj[0]) != 1 {
			return false
		}
		return g.Adj[0][0]
	}

	if len(g.Adj) == 2 {
		visited[0] = true
		visited[1] = true
		return g.Adj[0][1] && g.Adj[1][0] || g.Adj[0][0] || g.Adj[1][1]
	}

	visited[v] = true
	for child, b := range g.Adj[v] {
		if b {
			if child == v {
				return true
			}
			if visited[child] == false {
				if g.isCyclic(child, visited, v) == true {
					return true
				}
			} else if child != parent {
				return true
			}
		}
	}
	return false
}

func (g *Graph) isTree() bool {
	visited := make(map[int]bool)
	for i := range g.Adj {
		visited[i] = false
	}

	if g.isCyclic(0, visited, -1) == true {
		return false
	}

	for i := range g.Adj {
		if i > 0 && visited[i] == false {
			return false
		}
	}

	return true
}
