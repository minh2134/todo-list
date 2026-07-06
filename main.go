package main

import (
	"html/template"
	"net/http"
	"todo-list/internal/database"
	"todo-list/internal/task"
)

var (
	templ *template.Template
	db    database.Database
)

func main() {
	var err error
	templ, err = template.New("Task").ParseGlob("./templ/*.html")
	if err != nil {
		panic(err)
	}

	// TODO: read env var instead of hard code
	db, err = database.Open("todo.db")
	if err != nil {
		panic("Database error: " + err.Error())
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("GET /task", list)
	mux.HandleFunc("POST /task", add)
	mux.HandleFunc("DELETE /task", del)
	mux.HandleFunc("PATCH /task", edit)
	http.ListenAndServe(":3000", mux)
}

func list(w http.ResponseWriter, r *http.Request) {
	// TODO: handle query parameters to filter the result
	tsk := task.MakeTask("Task1", "hello")
	templ.ExecuteTemplate(w, "task", tsk)
}

func add(w http.ResponseWriter, r *http.Request) {
	tsk := task.MakeTask(r.PostFormValue("name"), r.PostFormValue("desc"))
	w.WriteHeader(204)
	db.InsertTask(tsk)
}

func del(w http.ResponseWriter, r *http.Request) {
	// TODO: delete the task (handle path and query params)
}

func edit(w http.ResponseWriter, r *http.Request) {
	// TODO: edit the task's data based on body (handle path)
}
