package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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
	mux.HandleFunc("DELETE /task/{id}", del)
	mux.HandleFunc("PATCH /task/{id}", edit)

	http.ListenAndServe(":3000", mux)
}

func list(w http.ResponseWriter, r *http.Request) {
	completed := database.ALL
	switch r.FormValue("completed") {
	case "true":
		completed = database.COMPLETED
	case "false":
		completed = database.INCOMPLETE
	}
	query := database.ListQuery{
		Name:      r.FormValue("name"),
		Completed: completed,
	}

	tsks, err := db.GetTasks(query)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
	}

	templ.ExecuteTemplate(w, "tasks", tsks)
}

func add(w http.ResponseWriter, r *http.Request) {
	tsk := task.MakeTask(r.PostFormValue("name"), r.PostFormValue("desc"))
	w.WriteHeader(204)
	id, err := db.InsertTask(tsk)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("inserted new task, id:", id)
}

func del(w http.ResponseWriter, r *http.Request) {
	// This always succeed even if the input is non-existent or non-sensical
	w.WriteHeader(204)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return
	}
	err = db.DeleteTask(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	// TODO: edit the task's data based on body (handle path)
}
