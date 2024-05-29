package handler

import (
	"todolist/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler { // конструктор для нового обработчика, накуя, пока не понял
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth") // Group - создаёт группу маршрутизаторов
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentify)
	{
		lists := api.Group("/lists") // зачем обертывать так всё ?
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllList)
			lists.GET("/id", h.getListById)
			lists.PUT("/id", h.updateList)
			lists.DELETE("/id", h.deleteList)

			items := lists.Group("id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItem)
				items.GET("/items_id", h.getItemById)
				items.PUT("/items_id", h.updateItem)
				items.DELETE("/items_id", h.deleteItem)
			}
		}
	}
	return router
}
