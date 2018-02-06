package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type cell struct {
	value    rune
	covered  bool
	backbone bool
	router   bool
	checked  bool
}

type file struct {
	cost        int
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
	start := time.Now()
	// Read File

	// file := readFile("final_round_2017.in/test.in")
	// file.grid[3][6].router = true
	// file.grid = addCoverage(file.grid, file.routerRange, 3, 6)
	// file.grid = connectToBackbone(file.grid, 3, 6)
	// file.grid[3][9].router = true
	// file.grid = addCoverage(file.grid, file.routerRange, 3, 9)
	// file.grid = connectToBackbone(file.grid, 3, 9)
	// printGrid(file.grid)
	// fmt.Println(calculateScore(file))
	//
	// newFile := readFile("final_round_2017.in/test.in")
	// newFile = run(newFile)
	// printGrid(newFile.grid)
	// fmt.Println(covered(newFile.grid))
	// fmt.Printf("cost: %v budget: %v\n", newFile.cost, newFile.budget)

	newFile := readFile(os.Args[1])
	newFile = run(newFile)
	printGrid(newFile.grid)
	fmt.Println(covered(newFile.grid))
	fmt.Printf("cost: %v budget: %v\n", newFile.cost, newFile.budget)

	// Place routers
	//		Detect collide with wall
	//		Update coverge
	// Connect routers to backbone
	// Calculate cost
	// Output file
	fmt.Printf("Execute time: %v", time.Since(start))
}

func run(newFile file) file {
	maxRating := (newFile.routerRange*2 + 1) * (newFile.routerRange*2 + 1)
	fmt.Println(maxRating)
	i := 0
	for {
		currentFile := newFile
		fmt.Printf("findNextRouters: %v\n", i)
		i++
		currentFile.grid, maxRating = findNextRouters(currentFile.grid, currentFile.routerRange, maxRating)
		currentFile.cost = calculateCost(currentFile.grid, currentFile.backbone, currentFile.router)
		if currentFile.cost > currentFile.budget {
			break
		}
		newFile = currentFile
		if fullyCovered(newFile.grid) {
			break
		}
	}
	return newFile
}

func calculateScore(model file) int {
	var score = 0
	score += 1000 * covered(model.grid)
	score += (model.budget - calculateCost(model.grid, model.backbone, model.router))
	score--
	return score
}

func calculateCost(grid [][]cell, backbone, router int) int {
	var cost = 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j].router {
				cost += router
			}
			if grid[i][j].backbone {
				cost += backbone
			}
		}
	}
	return cost
}

func connectToBackbone(grid [][]cell, x, y int) [][]cell {
	var x2, y2, bestDistance = 0, 0, math.MaxFloat64
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j].backbone {
				currentDistance := math.Sqrt(float64((x-i)*(x-i) + (y-j)*(y-j)))
				if currentDistance < bestDistance {
					x2 = i
					y2 = j
					bestDistance = currentDistance
				}
			}
		}
	}

	startX := x
	endX := x2
	if x > x2 {
		startX = x2
		endX = x
	}
	for i := startX; i <= endX; i++ {
		grid[i][y].backbone = true
	}

	startY := y
	endY := y2
	if y > y2 {
		startY = y2
		endY = y
	}
	for i := startY; i <= endY; i++ {
		grid[x2][i].backbone = true
	}
	return grid
}

func fullyCovered(grid [][]cell) bool {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j].value == target && !grid[i][j].checked {
				return false
			}
		}
	}
	return true
}

func findNextRouters(grid [][]cell, radius, maxRating int) ([][]cell, int) {
	var x, y, rating = 0, 0, 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if getRating(grid, radius, i, j) >= maxRating && maxRating != 0 {
				fmt.Printf("\tmaxRating: %v found\n", maxRating)
				fmt.Printf("\t\t x,y: %v %v\n", i, j)
				x = i
				y = j
				grid[x][y].router = true
				addCoverage(grid, radius, x, y)
				connectToBackbone(grid, x, y)
				rating = getRating(grid, radius, x, y)
			} else if getRating(grid, radius, i, j) > rating {
				x = i
				y = j
				rating = getRating(grid, radius, x, y)
			}

		}
	}
	grid[x][y].router = true
	addCoverage(grid, radius, x, y)
	connectToBackbone(grid, x, y)
	fmt.Printf("\trating: %v found\n", rating)
	fmt.Printf("\t\t x,y: %v %v\n", x, y)
	return grid, rating
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
	// fmt.Printf("Len grid: %v len grid[0]: %v\n", len(grid), len(grid[0]))
	for i := 0; i < (radius*2 + 1); i++ {
		for j := 0; j < (radius*2 + 1); j++ {
			if x-radius+i >= 0 && x-radius+i < len(grid) &&
				y-radius+j >= 0 && y-radius+j < len(grid[0]) {
				if grid[x-radius+i][y-radius+j].value != wall && grid[x-radius+i][y-radius+j].value != void {
					grid[x-radius+i][y-radius+j].checked = true
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

	// fmt.Printf("coords 1: %v %v\t coords 2: %v %v\n", x1, y1, x2, y2)

	if x1 > x2 {
		tmp := x1
		x1 = x2
		x2 = tmp
	}
	for i := x1; i < x2; i++ {
		if grid[i][y1].value == wall || grid[i][y2].value == wall {
			// fmt.Printf("\ti: %v y1: %v y2: %v %v\n", i, y1, y2, true)
			return true
		}
	}
	if y1 > y2 {
		tmp := y1
		y1 = y2
		y2 = tmp
	}
	for i := y1; i < y2; i++ {
		if grid[x1][i].value == wall || grid[x2][i].value == wall {
			// fmt.Printf("\ti: %v x1: %v x2: %v %v\n", i, x1, x2, true)
			return true
		}
	}
	// fmt.Printf("\t%v\n", false)
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
			} else if cell.checked {
				color.New(color.BgRed).Print(string(cell.value))
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

func getRating(grid [][]cell, radius, x, y int) int {
	rating := 0
	for i := 0; i < (radius*2 + 1); i++ {
		for j := 0; j < (radius*2 + 1); j++ {
			if x-radius+i >= 0 && x-radius+i < len(grid) &&
				y-radius+j >= 0 && y-radius+j < len(grid[0]) {
				if grid[x-radius+i][y-radius+j].value != wall && grid[x-radius+i][y-radius+j].value != void {
					if !isBlocked(grid, x, y, x-radius+i, y-radius+j) {
						if !grid[x-radius+i][y-radius+j].covered {
							rating++
						}
					}
				}
			}
		}
	}
	return rating
}
