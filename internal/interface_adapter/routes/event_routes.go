package routes

import (
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/interface_adapter/controller"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"example.com/EVENT-MANAGEMENT-SYSTEM/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func RegistereventsRoutes(routes *gin.Engine, eventController *controller.EventController, tokenRepo repository.TokenRepository) {
	authMiddleware := middlewares.AuthMiddleware(tokenRepo)

	eventGroup := routes.Group("/events")
	{
		// Protected routes (require valid authentication)
		eventGroup.Use(authMiddleware)
		{
			eventGroup.POST("", eventController.CreateEvent)
			eventGroup.PUT("/:id", eventController.UpdateEvent)
			eventGroup.DELETE("/:id", eventController.DeleteEvent)
			eventGroup.GET("/:id", eventController.GetEventByID)
			eventGroup.GET("", eventController.ListtAllEvents)
		}
	}
}
