package main

import (
	"database/sql"
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
	mux.HandleFunc("GET /task/{id}", editDiag)
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id not valid"))
		return
	}

	// database.ALL signifies a toggle signal here
	completed := database.ALL
	if r.FormValue("toggle") != "true" {
		switch r.FormValue("completed") {
		case "true":
			completed = database.COMPLETED
		case "false":
			completed = database.INCOMPLETE
		}
	} else {
		fmt.Println("we are toggling")
	}

	query := database.EditQuery{
		Id:        id,
		Name:      r.FormValue("name"),
		Desc:      r.FormValue("desc"),
		Completed: completed,
	}
	err = db.EditTask(query)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(os.Stderr, err)
	}
}

func editDiag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id not valid"))
		return
	}

	tsk, err := db.GetTask(id)
	if err == sql.ErrNoRows {
		w.WriteHeader(404)
		w.Write([]byte("id not found"))
	}

	data := struct {
		Id   int
		Task task.Task
	}{
		Id:   id,
		Task: tsk,
	}

	templ.ExecuteTemplate(w, "editTasks", data)
}
