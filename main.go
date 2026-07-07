package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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

	dbFile := os.Getenv("TODO_DB_FILE")
	if dbFile == "" {
		dbFile = "todo.db"
	}
	db, err = database.Open(dbFile)
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

	fmt.Println("Ready!")
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
