package main

import (
	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID `json: "id"`
	Name string    `json:"name"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:   dbUser.ID,
		Name: dbUser.Name,
	}
}
