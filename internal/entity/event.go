package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	StartTime   time.Time  `json:"starttime"`
	EndTime     time.Time  `json:"endtime"`
	Capacity    int        `json:"capacity"`
	IsPublic    bool       `json:"ispublic"`
	Status      string     `json:"status"`
	OrganizerID uuid.UUID  `json:"organizerid"`
	CreatedAt   time.Time  `json:"createdat"`
	UpdatedAt   time.Time  `json:"updatedat"`
	DeletedAt   *time.Time `json:"deletedat"`
}
