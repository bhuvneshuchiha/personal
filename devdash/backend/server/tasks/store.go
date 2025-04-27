package tasks

import (
)

type TaskStore interface {
	AddTask(t Task) error
	ListTask() ([]Task, error)
	UpdateTaskStatus(id int, status string) error
	DeleteTask(id int) error
}



// type InMemoryTask struct {
// 	tasks []Task
// }

// func (s *InMemoryTask) AddTask(t Task) error {
// 	s.tasks = append(s.tasks, t)
// 	return nil
// }
//
// func (s *InMemoryTask) ListTask() ([]Task, error) {
// 	// for _,v := range s.tasks {
// 	// 	fmt.Printf("Task Description :- ", v.Description)
// 	// }
// 	if len(s.tasks) > 0 {
// 		return s.tasks, nil
// 	}else {
// 		return nil, errors.New("No tasks found")
// 	}
// }
//
// func (s *InMemoryTask) UpdateTaskStatus(id int, status string) error {
// 	if id >=0 && len(s.tasks) > 0 {
// 		s.tasks[id].Status = status
// 		return nil
// 	}
// 	return errors.New("There was a problem updating the status of the task")
// }
//
// func (s *InMemoryTask) DeleteTask(id int) error {
// 	if id >= 0 && len(s.tasks) > 0 {
// 		s.tasks = slices.Delete(s.tasks, id, id + 1)
// 		return nil
// 	}
// 	return errors.New("Couldnt delete the task")
// }
//
