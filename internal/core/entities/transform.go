package entities

import (
	"github.com/google/uuid"
)

type ObjectPosition struct {
	ObjectID uuid.UUID `json:"object_id" db:"object_id"`
	X        float64   `json:"x" db:"x"`
	Y        float64   `json:"y" db:"y"`
	Z        float64   `json:"z" db:"z"`
}

type ObjectRotation struct {
	ObjectID uuid.UUID `json:"object_id" db:"object_id"`
	RX       float64   `json:"rx" db:"rx"`
	RY       float64   `json:"ry" db:"ry"`
	RZ       float64   `json:"rz" db:"rz"`
}

type ObjectScale struct {
	ObjectID uuid.UUID `json:"object_id" db:"object_id"`
	ScaleX   float64   `json:"scale_x" db:"scale_x"`
	ScaleY   float64   `json:"scale_y" db:"scale_y"`
	ScaleZ   float64   `json:"scale_z" db:"scale_z"`
}
