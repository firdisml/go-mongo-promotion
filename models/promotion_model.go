package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Promotion struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Category    string             `json:"category,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Shop        string             `json:"shop,omitempty" validate:"required"`
	State       string             `json:"state,omitempty" validate:"required"`
	Link        string             `json:"link,omitempty" validate:"required"`
}
