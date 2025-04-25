package repository

import (
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"github.com/gofrs/uuid"
)

type EventRepository interface {

	// CreateEvent creates a new event
	Create(event *entity.Event) error

	// Updatevent updates an existing event
	Update(event *entity.Event) error

	// Deleteevent deletes a event by its ID
	Delete(eventID uuid.UUID) error

	// Getevent returns a event by its ID
	GetByID(eventID uuid.UUID) (*entity.Event, error)

	// Getevents returns all event
	GetAll() ([]*entity.Event, error)
}
