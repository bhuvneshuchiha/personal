package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Position struct {
	Row int
	Col int
}

func CaptureKeypresses(loggerChan chan string) chan string {
	fmt.Println("Please start writing shit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		logger := scanner.Text()
		loggerChan <- logger
	}
	return loggerChan
}

func clearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}


func printGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(strings.Join(row, " "))
	}
}

func MoveSnake(g *Game, p *Position, pauseChan, quitChan chan bool, keylogger chan string) {
	grid := g.GameGrid
	gameSpeed := 300 * time.Millisecond

	direction := "right"

	for {
		select {
		case <-quitChan:
			fmt.Println("Quitting this boring ass game")
			return
		case key := <-keylogger:
			switch strings.ToLower(key) {
			case "w":
				direction = "up"
			case "s":
				direction = "down"
			case "a":
				direction = "left"
			case "d":
				direction = "right"
			case "q":
				close(quitChan)
				return
			}
		default:
			// Always move in the current direction
			grid[p.Row][p.Col] = " " // clear current position

			switch direction {
			case "up":
				p.Row--
			case "down":
				p.Row++
			case "left":
				p.Col--
			case "right":
				p.Col++
			}

			// Prevent out-of-bounds crash
			if p.Row < 0 {
				p.Row = 0
			}
			if p.Row >= len(grid) {
				p.Row = len(grid) - 1
			}
			if p.Col < 0 {
				p.Col = 0
			}
			if p.Col >= len(grid[0]) {
				p.Col = len(grid[0]) - 1
			}

			grid[p.Row][p.Col] = "*"

			printGrid(grid)

			time.Sleep(gameSpeed)
		}
	}
}

