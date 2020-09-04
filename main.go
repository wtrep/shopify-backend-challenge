package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func main() {
	bytes, err := bcrypt.GenerateFromPassword([]byte("adadadadaawdawdawdawdwadwadwadwa"), 10)
	if err != nil {
		panic(err)
	}
	print(bytes)
	//r := mux.NewRouter()
	//r.HandleFunc("/image", HandlePostImage).Methods("POST")
	//err := http.ListenAndServe(":8080", r)
	//if err != nil {
	//	panic(err)
	//}
}

func HandlePostImage(w http.ResponseWriter, r *http.Request) {

}
