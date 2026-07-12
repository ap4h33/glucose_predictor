package main

import (
	"github.com/ap4h33/glucose_predictor/internal/database"
)

type User struct {
	ID   int32  `json: "id"`
	Name string `json:"name"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:   dbUser.ID,
		Name: dbUser.Name,
	}
}
