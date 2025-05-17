package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhuvneshuchiha/devdash-backend/model"
	"github.com/bhuvneshuchiha/devdash-backend/server/tasks"
	_ "github.com/go-chi/chi/v5"
)

type Handler struct {
	DB *model.DatabaseObj
}

func (h *Handler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else {
		err := h.DB.AddTask(task)
		if err != nil {
			fmt.Println(err)
		} else {
			pretty, _ := json.Marshal(`{"message: data added successfully"}`)
			w.Write(pretty)
		}
	}
}

func (h *Handler) ListTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside check handler")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Inside the handler")
	getTasksFromDb, err := h.DB.GetTask();
	if err != nil {
		fmt.Println("Failed to retrieve the tasks")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	mp := make(map[string]any)
	mp["message"] = "Data fetched successfully"
	mp["tasks"] = getTasksFromDb
	pretty, e := json.MarshalIndent(getTasksFromDb, "","")
	if e != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Write(pretty)
}

func (h *Handler) UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
}
