package main

import (
	. "github.com/Garagator3000/cft-api/server"
	. "github.com/Garagator3000/cft-api/server/internal/handler"
	"github.com/Garagator3000/cft-api/server/internal/service"
)

func main() {
	service := service.NewService("/tmp/")
	handler := NewHandler(service)
	server := new(Server)

	server.Run("8080", handler.InitRoutes())
}
