// Revise closure and template literal string stuff
package main

import (
	"bytes"
	"fmt"
	"text/template"

)

type Player struct {
	Name string
	Score int
}

func (p Player) DisplayInfo() string {
	templ, err := template.New("test").Parse("Player {{.Name}} scored {{.Score}} points")
	if err != nil {
		return err.Error()
	}
	var buff bytes.Buffer
	er := templ.Execute(&buff, p)
	if er != nil {
		return er.Error()
	}
	return buff.String()
}

func (p *Player) AddPoints(points int) {
	p.Score += points
}

func ScoreBooster(factor int) func(int) int {
	return func(points int) int {
		points *= factor
		return points
	}
}

func ApplyBonus(player *Player, points int, bonusFunc func(int) int) {
	finalPoints := bonusFunc(points)
	player.Score += finalPoints
}

func main() {
	player1 := Player {
		Name: "Bhuvi",
		Score: 50,
	}
	info := player1.DisplayInfo()
	fmt.Println(info)

	p := Player{Name: "Bhuvnesh", Score: 20}
	fmt.Println(p.DisplayInfo())

	p.AddPoints(10)
	fmt.Println(p.DisplayInfo())

	double := ScoreBooster(2)
	ApplyBonus(&p, 5, double)
	fmt.Println(p.DisplayInfo())

}
