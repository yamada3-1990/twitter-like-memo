package main

import (
	"os"
	"twitter-like-memo/backend/app"
)

const (
	port         = "9000"
	imageDirPath = "images"
)

func main() {
	os.Exit(app.Server{
		Port:         port,
		ImageDirPath: imageDirPath,
	}.Run())
}
