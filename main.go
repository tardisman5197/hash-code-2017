package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Model of a cell
type cell struct {
	value    rune
	covered  bool
	backbone bool
	router   bool
	checked  bool
}

// Model of the input file
type file struct {
	cost         int
	routerRange  int
	backboneCost int
	routerCost   int
	budget       int
	grid         [][]cell
}

// Rune values of the expected cells
const (
	wall   = '#'
	target = '.'
	void   = '-'
)

// Funciton ran when program executed
func main() {
	start := time.Now()

	// Place routers
	//		Detect collide with wall
	//		Update coverge
	// Connect routers to backbone
	// Calculate cost
	// Output file

	newFile := readFile(os.Args[1])
	newFile = run(newFile)
	printGrid(newFile.grid)
	fmt.Printf("\nCovered: %v\n", covered(newFile.grid))
	fmt.Printf("Cost: %v Budget: %v\n", newFile.cost, newFile.budget)
	fmt.Printf("Execute time: %v\n", time.Since(start))
	fmt.Printf("Score: %v\n", calculateScore(newFile))
	writeFile(os.Args[2], newFile.grid)
}

// Places routers and backbone on the grid and calculates the cost
// Params: input file object
// Returns: output file object
func run(newFile file) file {
	// Init maxRating as maximum coverage a router can have
	maxRating := (newFile.routerRange*2 + 1) * (newFile.routerRange*2 + 1)
	fmt.Printf("maxRating: %v\n", maxRating)
	i := 0
	// Loop until budget exceeded or area fully covered
	for {
		fmt.Printf("findNextRouters: %v\n", i)
		i++
		// Get copy of file
		currentFile := newFile
		// Add routers that match maxRating, update maxRating with the biggest coverage
		currentFile.grid, maxRating = findNextRouters(currentFile.grid, currentFile.routerRange, maxRating)
		// Check if the budget has exceeded
		currentFile.cost = calculateCost(currentFile.grid, currentFile.backboneCost, currentFile.routerCost)
		if currentFile.cost > currentFile.budget {
			break
		}

		newFile = currentFile

		// Check if fully covered
		if fullyCovered(newFile.grid) {
			break
		}
	}
	return newFile
}

// Adds routers to the grid in positions depending on the coverage
// Params: grid
//				 radius of the routers
//				 max coverage of each routers
// Returns: output grid
//					largest coverage available in the grid
func findNextRouters(grid [][]cell, radius, maxRating int) ([][]cell, int) {
	// Init values
	var x, y, rating = 0, 0, 0
	// Loop though grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// if Cell rating is the max then place router
			// if max rating is 0 fully covered
			// else find highest rating available
			if getRating(grid, radius, i, j) >= maxRating && maxRating != 0 {
				fmt.Printf("\tmaxRating: %v found\n", maxRating)
				fmt.Printf("\t\t x,y: %v %v\n", i, j)
				// Add router to grid
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
	// Add router to grid with largest rating
	grid[x][y].router = true
	addCoverage(grid, radius, x, y)
	connectToBackbone(grid, x, y)
	fmt.Printf("\trating: %v found\n", rating)
	fmt.Printf("\t\t x,y: %v %v\n", x, y)
	return grid, rating
}

// Connects a point to the closest backbone cell
// Params: grid
//				 x coord of point to connect
//				 y coord of point to connect
// Returns: output grid
func worseConnectToBackbone(grid [][]cell, x, y int) [][]cell {
	// Init x2,y2 and set bestDistance to highest possible value
	var x2, y2, bestDistance = 0, 0, math.MaxFloat64
	// Loop though grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// Check if cell contains backbone
			if grid[i][j].backbone {
				// Calculate distance from point
				currentDistance := math.Sqrt(float64((x-i)*(x-i) + (y-j)*(y-j)))
				// Check if distance is the shortest
				if currentDistance < bestDistance {
					x2 = i
					y2 = j
					bestDistance = currentDistance
				}
			}
		}
	}

	// Set start and end x values
	startX := x
	endX := x2
	// if x is larger than x2 swap start and end x values
	if x > x2 {
		startX = x2
		endX = x
	}
	// Set grid from start to end x values to contain backbone
	for i := startX; i <= endX; i++ {
		grid[i][y].backbone = true
	}

	// Set start and end y values
	startY := y
	endY := y2
	// if y is larger than y2 swap start and end y values
	if y > y2 {
		startY = y2
		endY = y
	}
	// Set grid from start to end y values to contain backbone
	for i := startY; i <= endY; i++ {
		grid[x2][i].backbone = true
	}

	return grid
}

// Connects a point to the closest backbone cell
// Params: grid
//				 x coord of point to connect
//				 y coord of point to connect
// Returns: output grid
func connectToBackbone(grid [][]cell, x, y int) [][]cell {
	// Init x2,y2 and set bestDistance to highest possible value
	var x1, y1, bestDistance = 0, 0, math.MaxFloat64
	// Loop though grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// Check if cell contains backbone
			if grid[i][j].backbone {
				// Calculate distance from point
				currentDistance := math.Sqrt(float64((x-i)*(x-i) + (y-j)*(y-j)))
				// Check if distance is the shortest
				if currentDistance < bestDistance {
					x1 = i
					y1 = j
					bestDistance = currentDistance
				}
			}
		}
	}
	//
	// startX := x
	// endX := x2
	// if x > x2 {
	// 	startX = x2
	// 	endX = x
	// }
	// startY := y
	// endY := y2
	// if y > y2 {
	// 	startY = y2
	// 	endY = y
	// }
	//
	// dx := endX - startX
	// dy := endY - startY
	// de := math.Abs(float64(dy) / float64(dx))
	//
	// error := 0.0
	//
	// j := startY
	// for i := startX; i < endX; i++ {
	// 	grid[i][j].backbone = true
	// 	error = error + de
	// 	if error >= 0.5 {
	// 		j++
	// 		error -= 1.0
	// 	}
	// }

	dx := x1 - x
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		grid[x][y].backbone = true
		if x == x1 && y == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
	return grid
}

// Checks if the grid is fully covered
// Params: grid
// Returns: true if fully covered
func fullyCovered(grid [][]cell) bool {
	// Loop though grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// if cell is a target and is not covered, return false
			if grid[i][j].value == target && !grid[i][j].checked {
				return false
			}
		}
	}
	return true
}

// Adds the coverage of the router to the grid
// Params: grid
//				 radius of a router
//				 x value of the new router
//				 y value of the new router
// Returns: output grid
func addCoverage(grid [][]cell, radius, x, y int) [][]cell {
	// Loop through the possible cells for the new router
	for i := 0; i < (radius*2 + 1); i++ {
		for j := 0; j < (radius*2 + 1); j++ {
			// Check if Cell is not in the grid
			if x-radius+i >= 0 && x-radius+i < len(grid) &&
				y-radius+j >= 0 && y-radius+j < len(grid[0]) {
				// Check if Cell is not a wall or void Cell
				if grid[x-radius+i][y-radius+j].value != wall && grid[x-radius+i][y-radius+j].value != void {
					// Update Cell to be checked, for debug purposes
					grid[x-radius+i][y-radius+j].checked = true
					// Check Cell to router is not blocked
					if !isBlocked(grid, x, y, x-radius+i, y-radius+j) {
						// Update to Cell to be covered
						grid[x-radius+i][y-radius+j].covered = true
					}
				}
			}
		}
	}
	return grid
}

// Gets the amount of target cells covered
// Params: grid
// Returns: number of cells covered
func covered(grid [][]cell) int {
	// Init total
	var total = 0
	// Loop through grid
	for _, row := range grid {
		for _, cell := range row {
			// if cell covered inc total
			if cell.covered {
				total++
			}
		}
	}
	return total
}

// Gets the rating of a Cell based on the coverage if router was placed
// Params: grid
//				 radius of a router
//				 x,y value of Cell
// Returns: rating
func getRating(grid [][]cell, radius, x, y int) int {
	// Init rating
	rating := 0
	// Loop through the possible cells for the router
	for i := 0; i < (radius*2 + 1); i++ {
		for j := 0; j < (radius*2 + 1); j++ {
			// Check if Cell is not out of the grid
			if x-radius+i >= 0 && x-radius+i < len(grid) &&
				y-radius+j >= 0 && y-radius+j < len(grid[0]) {
				// Check if cell is not a wall or void
				if grid[x-radius+i][y-radius+j].value != wall && grid[x-radius+i][y-radius+j].value != void {
					// Check if path is blocked from router position
					if !isBlocked(grid, x, y, x-radius+i, y-radius+j) {
						// if Cell is not already covered, inc rating
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

// Check if path between two points is blocked by a wall
// Params: grid
//				 x,y values of first point
//				 x,y values of second point
// Returns: true if path blocked
func isBlocked(grid [][]cell, x1, y1, x2, y2 int) bool {
	// if x1 bigger than x2 swap x values
	if x1 > x2 {
		tmp := x1
		x1 = x2
		x2 = tmp
	}
	// Loop through x values between points
	for i := x1; i < x2; i++ {
		// Check if Cell is a wall at both point's y values
		if grid[i][y1].value == wall || grid[i][y2].value == wall {
			return true
		}
	}

	// if y1 bigger than y2 swap x values
	if y1 > y2 {
		tmp := y1
		y1 = y2
		y2 = tmp
	}
	// Loop through y values between points
	for i := y1; i < y2; i++ {
		// Check if Cell is a wall at both point's x values
		if grid[x1][i].value == wall || grid[x2][i].value == wall {
			return true
		}
	}
	return false
}

// Calculate the score of the grid, according to the Problem
// Params: output file object
// Returns: score
func calculateScore(model file) int {
	var score = 0
	// 1000 points per cell covered
	score += 1000 * covered(model.grid)
	// 1 point for each unit under budget
	score += (model.budget - calculateCost(model.grid, model.backboneCost, model.routerCost))
	return score
}

// Calculates the cost of the grid
// Params: grid
//				 cost of backbone
//				 cost of router
// Returns: cost of the grid
func calculateCost(grid [][]cell, backbone, router int) int {
	var cost = 0
	// Loop through cells in grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// If cell contains a router add to cost
			if grid[i][j].router {
				cost += router
			}
			// If cell contains backbone add to cost
			if grid[i][j].backbone {
				cost += backbone
			}
		}
	}
	// Remove cost of the inital backbone
	cost -= backbone
	return cost
}

// Prints the grid
// Params: grid
func printGrid(grid [][]cell) {
	//	- = void
	//	# = wall
	//	. = target
	//	blue = backbone
	//	green = covered
	//	yellow = router
	//	red = checked

	// colour importance:
	// 	router
	//	backbone
	//	covered
	//	checked
	fmt.Println()
	// Loop through grid
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
		}
		fmt.Println()
	}
}

// Reads input file
// Params: name of the input file
// Returns: the input file object
func readFile(filename string) file {
	// 	The first line contains the following numbers:
	// ​ ● H (1 ≤ H ≤ 1000)- the number of rows of the grid
	// ​ ● W (1 ≤ W ≤ 1000) - the number of columns of the grid
	// ​ ● R ( 1 ≤ R ≤ 10 ) - radius of a router range

	// The next line contains the following numbers:
	// ​ ● Pb (1 ≤ Pb ≤ 5) - price of connecting one cell to the backbone
	// ​ ● Pr (5 ≤ Pr ≤ 100) - price of one wireless router
	// ​ ● B ( 1 ≤ B ≤ 10 ) - maximum budget

	// The next line contains the following numbers:
	// ​ ● br , bc (0 ≤ br < H, 0 ≤ bc < W)- row and column of the initial cell that is already connected to the
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
			newFile.backboneCost, _ = strconv.Atoi(row[0])
			newFile.routerCost, _ = strconv.Atoi(row[1])
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

// Writes the output file
// Params: file name of the output file
//				 grid
func writeFile(filename string, grid [][]cell) {
	// The submission file must start with a line containing a single number N ( 0 ≤ N < W × H ) - the number of
	// cells connected to the backbone.

	// N next lines must specify the cells connected to the backbone, without repetitions and not including the
	// initial cell connected to the backbone that is specified in the problem statement. Each cell in the list must be
	// either neighbors with the initial backbone cell, or must appear in the list after one of its neighbors. Each line
	// in the list has to contain two numbers: r , c ( 0 ≤ r < H, 0 ≤ c < W)- respectively the row and the column of
	// each cell connected to the backbone.
	// The next line must contain a single number M ( 0 ≤ M ≤ W × H ) - the number of cells where routers are
	// placed.

	// M next lines must specify the cells where routers are placed without repetitions. Each of these lines must
	// contain two numbers: r , c ( 0 ≤ r < H, 0 ≤ c < W)- respectively the row and the column of each cell where a
	// router is placed
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	// Init values
	var backboneTotal, routerTotal, backbones, routers = 0, 0, "", ""

	// Loop though grid
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			// if cell has backbone, inc backboneTotal and add coords to list
			if grid[i][j].backbone {
				backboneTotal++
				backbones += strconv.Itoa(i) + " " + strconv.Itoa(j) + "\n"
			}
			// if cell has router, inc routerTotal and add coords to list
			if grid[i][j].router {
				routerTotal++
				routers += strconv.Itoa(i) + " " + strconv.Itoa(j) + "\n"
			}
		}
	}
	outputStr := strconv.Itoa(backboneTotal) + "\n"
	outputStr += backbones
	outputStr += strconv.Itoa(routerTotal) + "\n"
	outputStr += routers
	fmt.Fprintf(file, outputStr)
}
