package main

import (
	"github.com/wtrep/image-repo-backend/repo-backend"
)

func main() {
	// TODO check environment variables
	repo_backend.SetupAndServeRoutes()
}
