package Grid

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func InitGrid(rows, cols int) [][]string {
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
	}
	return grid
}

func PrintGrid(grid [][]string) {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			if grid[i][j] != "* " {
				fmt.Printf(". ")
			} else {
				fmt.Printf("* ")
			}
		}
		fmt.Println()
	}
}


func CaptureKeypresses(logger chan string) {
    reader := bufio.NewReader(os.Stdin)
    for {
        char, _, err := reader.ReadRune()
        if err != nil {
            log.Panic("System failed", err)
        } else {
            logger <- string(char)
        }
    }
}

func MoveSnake(grid [][]string, quitChan chan bool, keylogger chan string) {
    fmt.Println("Inside moveSnake")
    for {
        select {
        case key := <-keylogger:
            fmt.Println("Key received:", key)
            if key == "w" {
                fmt.Println("Going up")
                quitChan <- true
                return
            }
        }
    }
}
// func CaptureKeypresses(logger chan string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	line, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Panic("System failed", err)
// 	} else {
// 		logger <- line
// 	}
// }
//
// func MoveSnake(grid [][]string, quitChan chan bool, keylogger chan string) {
// 	fmt.Println("Inside moveSnake")
// 	for {
// 		select {
// 		case key := <-keylogger:
// 			fmt.Println("Key received:", key)
// 			if key == "w" {
// 				fmt.Println("Going up")
// 				quitChan <- true
// 				return
// 			}
// 		}
// 	}
// }
