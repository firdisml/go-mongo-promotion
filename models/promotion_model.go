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
	Created  time.Time          `json:"created" validate:"required"`
	Start    time.Time          `json:"start" validate:"required"`
	End      time.Time          `json:"end" validate:"required"`
	Visible  *bool              `json:"visible" validate:"required"`
}

type PromotionUpdate struct {
	Title    string    `json:"title,omitempty" validate:"required"`
	Category string    `json:"category,omitempty" validate:"required"`
	Shop     string    `json:"shop,omitempty" validate:"required"`
	State    string    `json:"state,omitempty" validate:"required"`
	Link     string    `json:"link,omitempty" validate:"required"`
	Start    time.Time `json:"start" validate:"required"`
	End      time.Time `json:"end" validate:"required"`
	Visible  *bool     `json:"visible,omitempty" validate:"required"`
	Picture  *bool     `json:"picture" validate:"required"`
}
