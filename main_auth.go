package main

import (
	"github.com/wtrep/image-repo-backend/auth"
)

func main() {
	// TODO check environment variables
	auth.SetupAndServeRoutes()
}
