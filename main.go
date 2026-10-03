package main

import (
	"os"

	"wallplay/internal/app"
)

func main() {
	os.Exit(app.New().Run())
}
