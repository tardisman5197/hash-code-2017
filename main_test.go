package main

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {
	result := readFile("final_round_2017.in/test.in")
	var expected file
	expected.routerRange = 3
	expected.backbone = 1
	expected.router = 100
	expected.budget = 220
	if result.routerRange != expected.routerRange {
		t.Fatalf("Read File Error: \n%v\n%v", result.routerRange, expected.routerRange)
	}
	if result.backbone != expected.backbone {
		t.Fatalf("Read File Error: \n%v\n%v", result.backbone, expected.backbone)
	}
	if result.router != expected.router {
		t.Fatalf("Read File Error: \n%v\n%v", result.router, expected.router)
	}
	if result.budget != expected.budget {
		t.Fatalf("Read File Error: \n%v\n%v", result.budget, expected.budget)
	}
	if result.grid[2][7].backbone != true {
		t.Fatalf("Read File Error: \n%v\n%v", result.grid[2][7].backbone, true)
	}
}

func TestCoverage(t *testing.T) {
	file := readFile("final_round_2017.in/test.in")
	file.grid[4][9].router = true
	file.grid = addCoverage(file.grid, file.routerRange, 4, 9)
	result := covered(file.grid)
	expected := 21
	if result != expected {
		t.Errorf("Coverage Error: %v != %v\nrouter: %v, %v", result, expected, 4, 9)
	}

	file = readFile("final_round_2017.in/test.in")
	file.grid[2][7].router = true
	file.grid = addCoverage(file.grid, file.routerRange, 2, 7)
	result = covered(file.grid)
	expected = 16
	if result != expected {
		t.Errorf("Coverage Error: %v != %v\nrouter: %v, %v", result, expected, 2, 7)
	}

	file = readFile("final_round_2017.in/test.in")
	file.grid[2][15].router = true
	file.grid = addCoverage(file.grid, file.routerRange, 2, 15)
	result = covered(file.grid)
	expected = 20
	if result != expected {
		t.Errorf("Coverage Error: %v != %v\nrouter: %v, %v", result, expected, 2, 15)
	}
}

func TestScore(t *testing.T) {
	file := readFile("final_round_2017.in/test.in")
	file.grid[3][6].router = true
	file.grid = addCoverage(file.grid, file.routerRange, 3, 6)
	file.grid = connectToBackbone(file.grid, 3, 6)
	file.grid[3][9].router = true
	file.grid = addCoverage(file.grid, file.routerRange, 3, 9)
	file.grid = connectToBackbone(file.grid, 3, 9)
	fmt.Println(calculateScore(file))
}
