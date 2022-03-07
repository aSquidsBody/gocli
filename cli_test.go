package gocli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCli2Graph(t *testing.T) {
	// test one child
	root := Command{Name: "root"}
	child1 := Command{Name: "c1"}

	cli := NewCli(&root)
	cli.AddChild(&root, &child1)

	adj := [][]bool{{false, true}, {false, false}}
	if !reflect.DeepEqual(cli2Graph(&cli).Adj, adj) {
		fmt.Println("AFTER")
		t.Errorf("cli2Graph failed for one child")
	}

	// test two children same level
	child2 := Command{Name: "c2"}
	cli.AddChild(&root, &child2)

	adj = [][]bool{{false, true, true}, {false, false, false}, {false, false, false}}
	if !reflect.DeepEqual(cli2Graph(&cli).Adj, adj) {
		t.Errorf("cli2Graph failed for two children")
	}

	// test two children cascading
	cli = NewCli(&root)
	cli.AddChild(&root, &child1)
	cli.AddChild(&child1, &child2)

	adj = [][]bool{{false, true, false}, {false, false, true}, {false, false, false}}
	if !reflect.DeepEqual(cli2Graph(&cli).Adj, adj) {
		t.Errorf("cli2Graph failed for two children cascading")
	}

	// test child with two sub-children
	// test two children cascading
	child3 := Command{Name: "c3"}
	cli = NewCli(&root)
	cli.AddChild(&root, &child1)
	cli.AddChild(&child1, &child2)
	cli.AddChild(&child1, &child3)

	adj = [][]bool{{false, true, false, false}, {false, false, true, true}, {false, false, false, false}, {false, false, false, false}}
	if !reflect.DeepEqual(cli2Graph(&cli).Adj, adj) {
		t.Errorf("cli2Graph failed for one child two sub-children")
	}

	// test cyclic one child
	cli = NewCli(&root)
	cli.AddChild(&root, &child1)
	cli.AddChild(&child1, &root)

	adj = [][]bool{{false, true}, {true, false}}
	if !reflect.DeepEqual(cli2Graph(&cli).Adj, adj) {
		t.Errorf("cli2Graph failed for one child cyclic")
	}
}

func TestIsCyclic(t *testing.T) {

	// test trivial graph
	g := Graph{
		Adj: [][]bool{},
	}
	if g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned true for trivial graph")
	}

	// test empty 2d graph
	g = Graph{
		Adj: [][]bool{{}},
	}
	if g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned true for empty 2d graph")
	}

	// test one element
	g = Graph{
		Adj: [][]bool{{false}},
	}
	if g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned true one element")
	}

	// test one element cyclic
	g = Graph{
		Adj: [][]bool{{true}},
	}
	if !g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned false one element cyclic")
	}

	// test two elements
	g = Graph{
		Adj: [][]bool{{false, false}, {false, false}},
	}
	if g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned true for two elements")
	}

	// test two elements cyclic
	g = Graph{
		Adj: [][]bool{{false, true}, {true, false}},
	}
	if !g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned false for two elements cyclic")
	}

	// test two elements identity cyclic
	g = Graph{
		Adj: [][]bool{{true, false}, {true, false}},
	}
	if !g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned false for two elements identity cyclic")
	}

	// test two elements not cyclic
	g = Graph{
		Adj: [][]bool{{false, true}, {false, false}},
	}
	if g.isCyclic(0, map[int]bool{}, -1) {
		t.Errorf("isCyclic returned true for two elements not cyclic")
	}

	// test three elements not cyclic
	g = Graph{
		Adj: [][]bool{{false, false, false}, {false, false, false}, {false, false, false}},
	}
	if g.isCyclic(0, map[int]bool{0: false, 1: false, 2: false}, -1) {
		t.Errorf("isCyclic returned true for three elements not cyclic")
	}

	// test three element chain not cyclic
	g = Graph{
		Adj: [][]bool{{false, true, false}, {false, false, true}, {false, false, false}},
	}
	if g.isCyclic(0, map[int]bool{0: false, 1: false, 2: false}, -1) {
		t.Errorf("isCyclic returned true for three element chain not cyclic")
	}

	// test three element chain cyclic
	g = Graph{
		Adj: [][]bool{{false, true, false}, {false, false, true}, {true, false, false}},
	}
	if !g.isCyclic(0, map[int]bool{0: false, 1: false, 2: false}, -1) {
		t.Errorf("isCyclic returned false for three element chain cyclic")
	}

	// test three element chain multiple cyclic
	g = Graph{
		Adj: [][]bool{{false, true, false}, {true, false, true}, {true, false, false}},
	}
	if !g.isCyclic(0, map[int]bool{0: false, 1: false, 2: false}, -1) {
		t.Errorf("isCyclic returned false for three element chain multiple cyclic")
	}

	// test three element chain identity cyclic
	g = Graph{
		Adj: [][]bool{{true, false, false}, {false, true, false}, {false, true, true}},
	}
	if !g.isCyclic(0, map[int]bool{0: false, 1: false, 2: false}, -1) {
		t.Errorf("isCyclic returned false for three element identity cyclic")
	}
}
