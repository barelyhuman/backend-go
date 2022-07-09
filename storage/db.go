package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

type TaskDetails struct {
	ID   int64  `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

func (t *TaskDetails) String() string {
	return "hello"
}

type DBClient struct {
	Client *sql.DB
}

func (s *DBClient) GetTasksAsJSON() (tasks []TaskDetails, err error) {
	res, err := s.Client.Query("select id,task,done from tasks")
	if err != nil {
		return
	}
	defer res.Close()
	for res.Next() {
		var (
			id   int64
			task string
			done bool
		)
		if err = res.Scan(&id, &task, &done); err != nil {
			return
		}
		tasks = append(tasks, TaskDetails{ID: id, Task: task, Done: done})
	}
	return
}

func (s *DBClient) TaskExists(id int64) bool {
	res, err := s.Client.Query("select id from tasks where id=$1", id)
	if err != nil {
		log.Println("failed", err)
		return false
	}
	defer res.Close()
	for res.Next() {
		if err = res.Scan(&id); err != nil {
			return false
		}
	}

	return true
}

func (s *DBClient) UpdateTask(task TaskDetails) (retId int, err error) {
	_, err = s.Client.Exec("update tasks set task = $2 , done = $3 where id=$1", task.ID, task.Task, task.Done)
	if err != nil {
		return
	}
	return
}

func (s *DBClient) InsertTask(task TaskDetails) (id int, err error) {
	_, err = s.Client.Exec("insert into tasks (task,done) values ($1,$2)", task.Task, task.Done)
	if err != nil {
		return
	}
	return
}

func (s *DBClient) RemoveTask(id int64) (err error) {
	_, err = s.Client.Exec("delete from tasks where id=$1;", id)
	if err != nil {
		return
	}
	return
}

func NewDB() *DBClient {
	connStr := os.Getenv("DATABASE_URL")
	if len(connStr) == 0 {
		panic("i need a database url to connect to.")
	}

	connStr, err := pq.ParseURL(connStr)

	if err != nil {
		panic("invalid connection string")
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("unable to connect to database, err: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("unable to connect to database, err: %v", err))
	}

	return &DBClient{
		Client: db,
	}
}
