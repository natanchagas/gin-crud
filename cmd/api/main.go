package main

import (
	"github.com/natanchagas/gin-crud/cmd/api/server"
	_ "github.com/natanchagas/gin-crud/config"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
