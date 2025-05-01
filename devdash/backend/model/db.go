package model

import (
	"database/sql"
	"fmt"

	"github.com/bhuvneshuchiha/devdash-backend/server/tasks"
	"github.com/jmoiron/sqlx"
)

type DatabaseObj struct {
	DB *sqlx.DB
}

func InitialiseDb(db *sqlx.DB) *DatabaseObj {
	return &DatabaseObj{DB: db}
}

func (d *DatabaseObj) AddTask(task tasks.Task) error {
	insertQuery := `
	INSERT INTO devdash (
	description,
	title,
	status
	)VALUES ($1, $2, $3);
	`
	_, err := d.DB.Exec(insertQuery, task.Description, task.Title, task.Status)
	if err != nil {
		fmt.Println("Couldnt save the task to the DB")
		return err
	} else {
		fmt.Println("Task successfully added")
	}
	return nil
}


func (d *DatabaseObj) GetTask() ([]tasks.Task, error){
	getQuery := `
	SELECT * FROM DEVDASH `

	rows, err := d.DB.Query(getQuery)

	if err != nil {
		fmt.Println("Couldnt get the tasks")
		return nil,err
	}
	defer rows.Close()

	var tasksList []tasks.Task
	for rows.Next(){
		var t tasks.Task
		err := rows.Scan(&t.ID, &t.Description, &t.Title, &t.Status)
		if err != nil {
			panic(err)
		}else{
			tasksList = append(tasksList, t)
		}
	}
	return tasksList, nil
}

func (d *DatabaseObj) RemoveTasks(id int) error {
	removeQuery := `
	DELETE FROM DEVDASH WHERE id=$1
	`
	tasks, err := d.DB.Query(removeQuery, id)
	if err != nil {
		fmt.Println("The task was not deleted")
		return err
	}else {
		fmt.Println(tasks)
	}
	return nil
}


func (d *DatabaseObj) GetTaskById(id int) (*sql.Rows, error) {
	getTaskById := `
	select * from devdash where id = $1`

	tasksById, err := d.DB.Query(getTaskById, id)
	if err != nil {
		fmt.Println("Couldnt get the task by ID")
		return nil, err
	}else {
		return tasksById, nil
	}
}











