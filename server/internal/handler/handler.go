package handler

import "github.com/Garagator3000/cft-api/server/internal/service"
import "github.com/gin-gonic/gin"

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	rest := router.Group("/file")
	{
		rest.GET("/all", handler.GetAll)
		rest.GET(":name", handler.Get)
		rest.POST("/new", handler.Post)
		rest.PUT("/update", handler.Put)
		rest.DELETE("/delete/:name", handler.Delete)
	}
	return router
}
