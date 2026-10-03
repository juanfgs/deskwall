package main

import (
	"os"

	"deskwall/internal/app"
)

func main() {
	os.Exit(app.New().Run())
}
