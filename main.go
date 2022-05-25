package main

import (
	"log"
	"net/http"

	"github.com/barelyhuman/tasks/debug"
	"github.com/barelyhuman/tasks/server"
	"github.com/barelyhuman/tasks/storage"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type UserAppState struct {
	Authenticated bool
}

type TaskReponse struct {
	Tasks    []storage.TaskDetails
	HasError bool
	Error    string
}

func HomeHandler(s *server.Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := &TaskReponse{}
	tasks, err := s.Storage.GetTasksAsJSON()

	if err != nil {
		response.HasError = true
		response.Error = err.Error()
		s.Views.Templates.ExecuteTemplate(w, "HomePage", response)
		return
	}

	response.Tasks = tasks
	response.HasError = false
	response.Error = ""

	s.Views.Templates.ExecuteTemplate(w, "HomePage", response)
}

func EditHandler(s *server.Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tasks, _ := s.Storage.GetTasks()

	err := s.Views.Templates.ExecuteTemplate(w, "EditPage", struct {
		Tasks string
	}{
		Tasks: tasks,
	})

	debug.DebugErr(err)
}

func UpdateTasksHandler(s *server.Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	tasks := []storage.TaskDetails{}

	debug.DebugLog(r.Form)

	for i, v := range r.Form["task"] {
		done := false
		if len(r.Form["done"]) > i && r.Form["done"][i] == "on" {
			done = true
		}
		tasks = append(tasks, storage.TaskDetails{
			Task: v,
			Done: done,
		})
	}

	s.Storage.SetTasks(tasks)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverRef := server.NewServer()
	serverRef.Router.GET("/", serverRef.CreateHandler(HomeHandler))
	serverRef.Router.GET("/edit", serverRef.CreateHandler(EditHandler))
	serverRef.Router.POST("/edit", serverRef.CreateHandler(UpdateTasksHandler))

	serverRef.StartServer()
}
