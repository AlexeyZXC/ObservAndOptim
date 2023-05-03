package models

import (
	"github.com/google/uuid"
)

type HasIDInterface interface {
	// GetID() uuid.UUID
	SetID(id string)
}

type User struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name"`
}

// func (u User) GetID() uuid.UUID {
// 	return u.ID
// }

func (u *User) SetID(id string) {
	//App.logger.Error("ssdfsdf")
	//todo: log the error
	u.ID, _ = uuid.FromBytes([]byte(id))
}

type Article struct {
	ID     uuid.UUID `json:"id,omitempty"`
	Title  string    `json:"title"`
	Text   string    `json:"string"`
	UserID uuid.UUID `json:"user_id"`
}

// func (a Article) GetID() uuid.UUID {
// 	return a.ID
// }

func (a *Article) SetID(id string) {
	//todo: log the error
	a.ID, _ = uuid.FromBytes([]byte(id))
}
