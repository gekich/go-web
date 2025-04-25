package main

import (
	"github.com/gekich/go-web/internal/server"
)

// @title Swagger Example API
// @version 1.0
// @description Go-web swagger example.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:3000
// @BasePath /
func main() {
	s := server.New()
	s.Init()
	s.Run()
}
