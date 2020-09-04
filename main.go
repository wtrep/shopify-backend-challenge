package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/image", HandlePostImage).Methods("POST")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func HandlePostImage(w http.ResponseWriter, r *http.Request) {

}
