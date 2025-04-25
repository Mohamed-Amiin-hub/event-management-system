package controller

import (
	"net/http"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// User Controller handles event requests
type EventController struct {
	eventService service.EventService
}

// NewEventController creates a new EventController instance
func NewEventController(eventService service.EventService) *EventController {
	return &EventController{eventService: eventService}

}

// CreateEvent handles the creation of a new event
func (c *EventController) CreateEvent(ctx *gin.Context) {
	var event entity.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	OrganizerID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "organizer ID is required"})
		return
	}

	createdEvent, err := c.eventService.CreateEvent(
		event.Title,
		event.Description,
		event.Location,
		event.Capacity,
		event.Status,
		OrganizerID.(uuid.UUID),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, createdEvent)
}

// Updateevent handles the update of an existing event
func (c *EventController) UpdateEvent(ctx *gin.Context) {
	var event entity.Event

	eventIdparam := ctx.Param("id")
	eventID, err := uuid.FromString(eventIdparam)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid event ID"})
		return
	}

	// Bind JSON input to event struct
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = eventID

	// Call service to update event
	if err := c.eventService.UpdateEvent(&event); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

func (c *EventController) DeleteEvent(ctx *gin.Context) {

	eventIdparam := ctx.Param("id")
	eventID, err := uuid.FromString(eventIdparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := c.eventService.DeleteEvent(eventID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}

// GeteventByID handles retrieving a event by its ID
func (c *EventController) GetEventByID(ctx *gin.Context) {
	eventIdparam := ctx.Param("id")
	eventID, err := uuid.FromString(eventIdparam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := c.eventService.GetEventByID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (c *EventController) ListtAllEvents(ctx *gin.Context) {
	events, err := c.eventService.ListEvent()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
