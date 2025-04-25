package gateway

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
	"github.com/gofrs/uuid"
)

type EventRepositoryimpl struct {
	db *sql.DB
}

// Create implements repository.EventRepository.
func (e *EventRepositoryimpl) Create(event *entity.Event) error {

	log.Printf("Inserting into events: %+v", event)

	query := `INSERT INTO events (
		id, title, description, location, start_time, end_time,
		capacity, is_public, status, organizer_id, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11, $12
	)`

	result, err := e.db.Exec(query,
		event.ID, event.Title, event.Description, event.Location,
		event.StartTime, event.EndTime, event.Capacity, event.IsPublic,
		event.Status, event.OrganizerID, event.CreatedAt, event.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error inserting event: %v\nQuery: %s", err, query)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	log.Printf("Rows affected: %d", rowsAffected)
	return nil
}

// Delete implements repository.EventRepository.
func (e *EventRepositoryimpl) Delete(eventID uuid.UUID) error {

	query := `DELETE FROM events WHERE id = $1`
	result, err := e.db.Exec(query, eventID)
	if err != nil {
		log.Printf("Error deleting event with ID %v: %v", eventID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No event found with ID: %v", eventID)
		return nil
	}

	log.Printf("Deleted event with ID: %v", eventID)
	return nil
}

// GetAll implements repository.EventRepository.
func (e *EventRepositoryimpl) GetAll() ([]*entity.Event, error) {

	var events []*entity.Event

	query := `SELECT id, title, description, location, start_time, end_time,
		capacity, is_public, status, organizer_id, created_at, updated_at, deleted_at
		FROM events`

	rows, err := e.db.Query(query)
	if err != nil {
		log.Printf("Error retrieving events: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event entity.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.Location,
			&event.StartTime, &event.EndTime, &event.Capacity, &event.IsPublic,
			&event.Status, &event.OrganizerID, &event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating events: %v", err)
		return nil, err
	}

	return events, nil
}

// GetdByID implements repository.EventRepository.
func (e *EventRepositoryimpl) GetByID(eventID uuid.UUID) (*entity.Event, error) {
	var event entity.Event

	query := `SELECT id, title, description, location, start_time, end_time,
		capacity, is_public, status, organizer_id, created_at, updated_at, deleted_at
		FROM events WHERE id = $1`

	err := e.db.QueryRow(query, eventID).Scan(
		&event.ID, &event.Title, &event.Description, &event.Location,
		&event.StartTime, &event.EndTime, &event.Capacity, &event.IsPublic,
		&event.Status, &event.OrganizerID, &event.CreatedAt, &event.UpdatedAt, &event.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No event found with ID: %v", eventID)
			return nil, fmt.Errorf("event not found")
		}
		log.Printf("Error retrieving event by ID: %v", err)
		return nil, err
	}

	return &event, nil
}

// Update implements repository.EventRepository.
func (e *EventRepositoryimpl) Update(event *entity.Event) error {
	query := `UPDATE events
		SET title = $2, description = $3, location = $4, start_time = $5,
			end_time = $6, capacity = $7, is_public = $8, status = $9,
			organizer_id = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	result, err := e.db.Exec(query,
		event.ID, event.Title, event.Description, event.Location,
		event.StartTime, event.EndTime, event.Capacity, event.IsPublic,
		event.Status, event.OrganizerID,
	)
	if err != nil {
		log.Printf("Error updating event with ID %v: %v", event.ID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No event found with ID: %v", event.ID)
		return nil
	}

	log.Printf("Event updated successfully: %v", event.ID)
	return nil
}

// factory function to create an instance of EventRepository
func NewEventRepository(db *sql.DB) repository.EventRepository {
	return &EventRepositoryimpl{db: db}
}
