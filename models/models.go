package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username"      bson:"username"`
	Password  string             `json:"password"      bson:"password"`
	Firstname string             `json:"firstname"`
	Lastname  string             `json:"lastname"`
	Age       int                `json:"age"`
}

type Admin struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username"      bson:"username"`
	Password string             `json:"password"      bson:"password"`
}

type News struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title"`
	Body       string             `json:"body"`
	TimeOfCast time.Time          `json:"timeofcast"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
