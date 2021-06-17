package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./storage"))))
	http.HandleFunc("/static/avatar", handleFileUpload)

	println("Server listening at http://localhost:8000")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad request", 400)
		return
	}

	file, _, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	defer file.Close()

	f, err := os.Create("foo.png")
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("Created\n"))
}
