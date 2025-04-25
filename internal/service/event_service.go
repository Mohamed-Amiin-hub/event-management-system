package service

import (
	"fmt"
	"log"
	"time"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"github.com/gofrs/uuid"
)

type EventService interface {
	CreateEvent(title, description, location string, capacity int, status string, OrganizerID uuid.UUID) (*entity.Event, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(eventID uuid.UUID) error
	GetEventByID(eventID uuid.UUID) (*entity.Event, error)
	ListEvent() ([]*entity.Event, error)
}

// userServiceImpl is the implementation of UserService.
type EventServiceImpl struct {
	repo      repository.EventRepository
	tokenRepo repository.TokenRepository
}

// CreateEvent implements eventService.
func (s *EventServiceImpl) CreateEvent(title string, description string, location string, capacity int, status string, OrganizerID uuid.UUID) (*entity.Event, error) {
	// Generate a new UUID for the event ID
	neoEvent, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	// Create a new event instance
	newEvent := &entity.Event{
		ID:          neoEvent,
		Title:       title,
		Description: description,
		Location:    location,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		Capacity:    capacity,
		IsPublic:    true,
		Status:      status,
		OrganizerID: OrganizerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Log the new event creation attempt
	log.Printf("create Event: %+v", newEvent)

	// Save the new event to the repository
	err = s.repo.Create(newEvent)
	if err != nil {
		log.Printf("failed to create event: %v", err)
	}
	return newEvent, nil

}

// DeleteEvent implements eventService.
func (s *EventServiceImpl) DeleteEvent(eventID uuid.UUID) error {
	_, err := s.repo.GetByID(eventID)
	if err != nil {
		return fmt.Errorf("could not find event with ID %s: %v", eventID, err)
	}

	if err := s.repo.Delete(eventID); err != nil {
		return fmt.Errorf("failed to delete event with ID %s: %v", eventID, err)
	}

	log.Printf("Successfully deleted event with ID %s", eventID)
	return nil
}

// GetEventByID implements eventService.
func (s *EventServiceImpl) GetEventByID(eventID uuid.UUID) (*entity.Event, error) {
	event, err := s.repo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event with ID %s: %v", eventID, err)
	}

	return event, nil
}

// ListEvent implements eventService.
func (s *EventServiceImpl) ListEvent() ([]*entity.Event, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all event: %v", err)
	}

	return events, nil
}

// UpdateEvent implements eventService.
func (s *EventServiceImpl) UpdateEvent(event *entity.Event) error {
	_, err := s.repo.GetByID(event.ID)
	if err != nil {
		return fmt.Errorf("could not find event with ID %s", event.ID)
	}

	if err := s.repo.Update(event); err != nil {
		return fmt.Errorf("failed to update event with ID %s: %v", event.ID, err)
	}

	return nil
}

func NewEventService(eventRepo repository.EventRepository, tokenRepo repository.TokenRepository) EventService {
	return &EventServiceImpl{
		repo:      eventRepo,
		tokenRepo: tokenRepo,
	}
}
