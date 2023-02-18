package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promotion struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Title    string             `json:"title,omitempty" validate:"required"`
	Category string             `json:"category,omitempty" validate:"required"`
	Shop     string             `json:"shop,omitempty" validate:"required"`
	State    string             `json:"state,omitempty" validate:"required"`
	Link     string             `json:"link,omitempty" validate:"required"`
	Created  primitive.DateTime `json:"created"`
	Start    time.Time          `json:"start" validate:"required"`
	End      time.Time          `json:"end" validate:"required"`
	Visible  bool               `json:"visible"`
}
