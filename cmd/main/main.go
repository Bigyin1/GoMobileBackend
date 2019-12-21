package main

import (
	"github.com/Bigyin1/GoMobileBackend/app"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	application := app.InitApp()
	application.StartApp()
}
