package main

import (
	"github.com/hrshadhin/fiber-go-boilerplate/cmd/server"
	_ "github.com/hrshadhin/fiber-go-boilerplate/docs" // load API Docs files (Swagger)
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
)

func main() {

	// setup various configuration for app
	config.LoadAllConfigs(".env")

	server.Serve()
}
