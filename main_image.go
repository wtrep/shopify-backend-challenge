package main

import (
	"github.com/wtrep/image-repo-backend/image"
)

func main() {
	// TODO check environment variables
	image.SetupAndServeRoutes()
}
