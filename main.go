package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("GET /clicked", clicked)
	http.ListenAndServe(":3000", mux)
}

func clicked(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You clicked!"))
}
