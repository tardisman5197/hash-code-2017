package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type cell struct {
	value    rune
	covered  bool
	backbone bool
	router   bool
}

type file struct {
	routerRange int
	backbone    int
	router      int
	budget      int
	grid        [][]cell
}

const (
	wall   = '#'
	target = '.'
	void   = '-'
)

func main() {
	// Read File
	newFile := readFile("final_round_2017.in/test.in")
	newFile.grid[2][7].router = true
	newFile.grid = addCoverage(newFile.grid, newFile.routerRange, 2, 7)
	printGrid(newFile.grid)

	// Place routers
	//		Detect collide with wall
	//		Update coverge
	// Connect routers to backbone
	// Calculate cost
	// Output file
}

func readFile(filename string) file {
	fmt.Printf("Read File: %v\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	index := 0
	var newFile file
	// Loop through the file
	for s.Scan() {
		row := strings.Fields(s.Text())
		// fmt.Printf("%v : ", index)
		// fmt.Println(row)
		switch index {
		case 0:
			rowNum, _ := strconv.Atoi(row[0])
			colNum, _ := strconv.Atoi(row[1])
			grid := make([][]cell, rowNum)
			for i := 0; i < rowNum; i++ {
				grid[i] = make([]cell, colNum)
			}
			newFile.grid = grid
			newFile.routerRange, _ = strconv.Atoi(row[2])
		case 1:
			newFile.backbone, _ = strconv.Atoi(row[0])
			newFile.router, _ = strconv.Atoi(row[1])
			newFile.budget, _ = strconv.Atoi(row[2])
		case 2:
			x, _ := strconv.Atoi(row[0])
			y, _ := strconv.Atoi(row[1])
			newFile.grid[x][y].backbone = true
		default:
			for i, char := range row[0] {
				newFile.grid[index-3][i].value = char
			}
		}
		index++
	}
	return newFile
}

func addCoverage(grid [][]cell, radius, x, y int) [][]cell {
	fmt.Printf("Len grid: %v len grid[0]: %v", len(grid), len(grid[0]))
	for i := 0; i < (radius*2 + 1); i++ {
		for j := 0; j < (radius*2 + 1); j++ {
			if x-radius+i >= 0 && x-radius+i < len(grid) &&
				y-radius+j >= 0 && y-radius+j < len(grid[0]) {
				if grid[x-radius+i][y-radius+j].value != wall && grid[x-radius+i][y-radius+j].value != void {
					if !isBlocked(grid, x, y, x-radius+i, y-radius+j) {
						grid[x-radius+i][y-radius+j].covered = true
					}
				}
			}
		}
	}
	return grid
}

func isBlocked(grid [][]cell, x1, y1, x2, y2 int) bool {
	distX := x1 - x2
	distY := y1 - y2
	right := false
	down := false
	if distX < 0 {
		right = true
		distX = -distX
	}
	if distY < 0 {
		down = true
		distY = -distY
	}
	// fmt.Printf("distX: %v\tdistY: %v\tright: %v\tdown: %v\n", distX, distY, right, down)
	// fmt.Printf("coords 1: %v %v\t coords 2: %v %v\n", x1, y1, x2, y2)
	for i := 0; i < distX; i++ {
		var cellValue rune
		if right {
			cellValue = grid[x1-i][y1].value
		} else {
			cellValue = grid[x1+i][y1].value
		}
		if cellValue == wall {
			return true
		}
		if right {
			cellValue = grid[x1-i][y2].value
		} else {
			cellValue = grid[x1+i][y2].value
		}
		if cellValue == wall {
			return true
		}
	}
	for i := 0; i < distY; i++ {
		var cellValue rune
		if down {
			cellValue = grid[x1][y1-i].value
		} else {
			cellValue = grid[x1][y1+i].value
		}
		if cellValue == wall {
			return true
		}
		if down {
			cellValue = grid[x2][y1-i].value
		} else {
			cellValue = grid[x2][y1+i].value
		}
		if cellValue == wall {
			return true
		}
	}

	return false
}

func printGrid(grid [][]cell) {
	fmt.Println()
	for _, row := range grid {
		for _, cell := range row {
			if cell.router {
				color.New(color.BgYellow).Print(string(cell.value))
			} else if cell.backbone {
				color.New(color.BgCyan).Print(string(cell.value))
			} else if cell.covered {
				color.New(color.BgGreen).Print(string(cell.value))
			} else {
				fmt.Print(string(cell.value))
			}
			//fmt.Print(string(cell.value))
		}
		fmt.Println()
	}
}

func covered(grid [][]cell) int {
	var total = 0
	for _, row := range grid {
		for _, cell := range row {
			if cell.covered {
				total++
			}
		}
	}
	return total
}
