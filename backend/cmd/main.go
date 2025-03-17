package main

import (
	"mercari-build-training/app"
	"os"
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
