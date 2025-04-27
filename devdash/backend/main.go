package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/bhuvneshuchiha/devdash-backend/api"
	"github.com/bhuvneshuchiha/devdash-backend/model"
)

func ConnectDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=bhuvnesh dbname=devdash sslmode=disable")
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}


func main() {

	db := ConnectDatabase()
	db_obj := model.InitialiseDb(db)
	r := chi.NewRouter()
	h := &api.Handler{
		DB: db_obj,
	}

	r.Post("/addTask", h.AddTaskHandler)
	r.Get("/getTask", h.ListTaskHandler)
	r.Post("/updateTask", h.UpdateTaskStatusHandler)
	r.Post("/deleteTask/{id}", h.DeleteTaskHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
