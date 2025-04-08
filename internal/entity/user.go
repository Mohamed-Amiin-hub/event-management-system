package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

// User represents a user entity
type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	IsActive  bool       `json:"is_active"`
	LastLogin time.Time  `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
