package storage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type TaskDetails struct {
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type DBClient struct {
	*redis.Client
}

var ctx = context.Background()

// TODO: change to options pattern
func (s *DBClient) GetTasks() (tasks string, err error) {
	tasks, err = s.Get(ctx, "tasks").Result()
	return
}

func (s *DBClient) GetTasksAsJSON() (tasks []TaskDetails, err error) {
	tasksAsString, err := s.Get(ctx, "tasks").Result()
	json.Unmarshal([]byte(tasksAsString), &tasks)
	return
}

func (s *DBClient) SetTasks(tasks []TaskDetails) (err error) {
	tasksAsBytes, err := json.Marshal(tasks)
	if err != nil {
		return
	}
	s.Set(ctx, "tasks", string(tasksAsBytes), time.Duration(time.Hour*24)).Result()
	return
}

func NewDB() *DBClient {
	opts, _ := redis.ParseURL(os.Getenv("REDIS_CONNECTION_STRING"))

	rdb := redis.NewClient(opts)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panic("Failed to connect to client", err)
	}

	return &DBClient{
		Client: rdb,
	}
}
