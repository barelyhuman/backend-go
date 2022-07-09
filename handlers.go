package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/barelyhuman/tasks/debug"
	"github.com/barelyhuman/tasks/server"
	"github.com/barelyhuman/tasks/storage"
	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
)

type UserAppState struct {
	Authenticated bool
}

type TaskResponse struct {
	Tasks    []storage.TaskDetails
	HasError bool
	Error    string
}

type TasksState struct {
	storage.TaskDetails
	delete bool
}

func HomeHandler(s *server.Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := &TaskResponse{}
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
	tasks, _ := s.Storage.GetTasksAsJSON()
	response := struct {
		Tasks     string
		CSRFField template.HTML
	}{}

	// if err != nil {
	// 	// response.HasError = true
	// 	// response.Error = err.Error()
	// 	s.Views.Templates.ExecuteTemplate(w, "EditPage", response)
	// 	return
	// }

	str, _ := json.Marshal(tasks)
	response.Tasks = string(str)
	response.CSRFField = csrf.TemplateField(r)
	// response.HasError = false
	// response.Error = ""

	err := s.Views.Templates.ExecuteTemplate(w, "EditPage", response)

	debug.DebugErr(err)
}

func UpdateTasksHandler(s *server.Server, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	debug.DebugLog(r.Form)

	tasks := []TasksState{}

	for i, formTask := range r.Form["task"] {
		formId := r.Form["id"][i]
		idInt, _ := strconv.Atoi(formId)
		id := int64(idInt)
		basePayload := TasksState{}
		basePayload.Task = formTask
		basePayload.ID = id
		tasks = append(tasks, basePayload)
	}

	for _, idsToMark := range r.Form["done"] {
		idInt, _ := strconv.Atoi(idsToMark)
		id := int64(idInt)
		for i, ts := range tasks {
			if ts.ID == id {
				tasks[i].Done = true
				break
			}
		}
	}

	for _, idsToDel := range r.Form["delete"] {
		idInt, _ := strconv.Atoi(idsToDel)
		id := int64(idInt)
		for i, ts := range tasks {
			if ts.ID == id {
				tasks[i].delete = true
				break
			}
		}
	}

	for _, taskItem := range tasks {
		if taskItem.delete {
			s.Storage.RemoveTask(taskItem.ID)
			continue
		}
		if taskItem.ID != 0 && s.Storage.TaskExists(taskItem.ID) {
			s.Storage.UpdateTask(taskItem.TaskDetails)
		} else {
			s.Storage.InsertTask(taskItem.TaskDetails)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
