
package main

import (
	"github.com/bhuvneshuchiha/snake_game/game"
)

func main() {

	snaky := game.Game{
		IsRunning:true,
		IsPaused: false,
		GameGrid: game.InitGrid(10, 10),
	}

	position := game.Position{
		Row : 0,
		Col : 0,
	}

	quitChan := make(chan bool)
	pauseChan := make(chan bool, 1)
	loggerChan := make(chan string, 10)

	go game.CaptureKeypresses(loggerChan)
	go snaky.ManageGameState(quitChan, pauseChan)
	go game.MoveSnake(&snaky, &position, pauseChan,quitChan, loggerChan)

	<- quitChan

}
