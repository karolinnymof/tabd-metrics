package models

import "github.com/google/uuid"

type Rating struct {
	ID      uuid.UUID `json:"id"`
	BookID  uuid.UUID `json:"book_id"`
	UserID  uuid.UUID `json:"user_id"`
	Note    int       `faker:"num" json:"note"`
	Comment string    `faker:"sentence" json:"comment"`
}
