package entities

import (
	"github.com/google/uuid"
)

type Object struct {
	ID          uuid.UUID `json:"id" db:"object_id"`
	Name        string    `json:"name" db:"object_name"`
	Description string    `json:"description,omitempty" db:"object_description"`
	File        []byte    `json:"-" db:"object_file"`
	OwnerID     int       `json:"owner_id" db:"owner"`
	Movable     bool      `json:"movable" db:"movable"`
	Printable   bool      `json:"printable" db:"printable"`
}
