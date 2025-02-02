package main

import (
	"github.com/hafidzhz/ihsansolusi-test/cmd/server"
	_ "github.com/hafidzhz/ihsansolusi-test/docs" // load API Docs files (Swagger)
	"github.com/hafidzhz/ihsansolusi-test/pkg/config"
)

func main() {
	config.LoadAllConfigs(".env")

	server.Serve()
}
