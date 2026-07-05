package main

import (
	"html/template"
	"net/http"
	"todo-list/internal/task"
)

var templ *template.Template

func main() {
	var err error
	templ, err = template.New("Task").ParseGlob("./templ/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /clicked", clicked)

	mux.HandleFunc("GET /task", list)
	mux.HandleFunc("POST /task", add)
	mux.HandleFunc("DELETE /task", del)
	mux.HandleFunc("PATCH /task", edit)
	http.ListenAndServe(":3000", mux)
}

func clicked(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You clicked!"))
}

func list(w http.ResponseWriter, r *http.Request) {
	// TODO: handle query parameters to filter the result
	tsk := task.MakeTask("Task1", "hello")
	templ.ExecuteTemplate(w, "task", tsk)
}

func add(w http.ResponseWriter, r *http.Request) {
	// TODO: add the posted task
}

func del(w http.ResponseWriter, r *http.Request) {
	// TODO: delete the task (handle path and query params)
}

func edit(w http.ResponseWriter, r *http.Request) {
	// TODO: edit the task's data based on body (handle path)
}
