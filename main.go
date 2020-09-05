package main

import (
	"github.com/wtrep/image-repo-backend/backend"
)

func main() {
	// TODO check environment variables
	backend.SetupAndServeRoutes()
}
