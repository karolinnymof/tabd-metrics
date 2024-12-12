package models

import "github.com/google/uuid"

type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `faker:"word" json:"title"`
	Author      string    `faker:"word" json:"author"`
	Genre       string    `faker:"word" json:"genre"`
	Description string    `faker:"sentence" json:"description"`
}
