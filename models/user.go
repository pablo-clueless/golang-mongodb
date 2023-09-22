package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username" binding:"required"`
	Email     string    `json:"email" bson:"email" binding:"required"`
	Password  string    `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
