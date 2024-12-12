package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `faker:"username" json:"name"`
	Email       string    `faker:"email" json:"email"`
	Preferences string    `faker:"word" json:"preferences"`
}
