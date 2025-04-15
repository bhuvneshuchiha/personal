package main

import (
	"fmt"

	"github.com/bhuvneshuchiha/devdash-backend/internal/tasks"
)

func main() {
	fmt.Println("Welcome to devdash bitches")
	store := tasks.InMemoryTask{}
	store.AddTask(tasks.Task{
		ID: 1,
		Title: "Play Sports",
		Description: "Play football",
		Status: "Todo",
	})

	tasksToShow, err := store.ListTask()
	if err != nil {
		fmt.Println("Doomed")
	}else {
		for _, v := range tasksToShow {
			fmt.Printf("ID :- %d\nTitle :- %s\nDescription :- %s\nStatus :- %s\n",
				v.ID, v.Title, v.Description, v.Status)
		}
	}
}
