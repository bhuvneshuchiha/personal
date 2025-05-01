package main

import Grid "github.com/bhuvneshuchiha/snake2/grid"

func main() {
	rows := 20
	cols := 40
	grid := Grid.InitGrid(rows, cols)
	grid[0][0] = "* "
	Grid.PrintGrid(grid)

	logger := make(chan string)
	quit := make(chan bool)

	go Grid.CaptureKeypresses(logger)
	go Grid.MoveSnake(grid, quit, logger)

	<-quit

}
