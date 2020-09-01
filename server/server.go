package main

import (
	"net/http"
	"github.com/google/uuid"
)

type Package struct {
	Owner string
	Name  string
	UUID  uuid.UUID
}

func handleGetPackages(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/packages", handleGetPackages)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
