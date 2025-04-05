package game

import (
	"fmt"
	"time"
)

type Game struct {
	IsRunning  bool
	IsPaused   bool
	GameGrid [][]string
}

func InitGrid(rows, cols int) [][]string {
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
	}
	return grid
}

func (g *Game) ManageGameState(quitChan, pauseChan chan bool) {
	for {
		fmt.Println("Game is running")
		select {
		case <-quitChan:
			fmt.Println("Game over!!!")
			return
		case <-pauseChan:
			g.IsPaused = true
			fmt.Println("Game is paused for 10 seconds")
			for i := 0; i < 100; i++ {
				select {
				case <-quitChan:
					fmt.Println("quitting the game!!, exiting the paused state")
					return
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
			g.IsPaused = false
			fmt.Println("Game resumed")
		}
	}
}
