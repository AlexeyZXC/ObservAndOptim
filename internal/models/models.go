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
type UserAsArray struct {
	ID   [16]byte `json:"id,omitempty"`
	Name string   `json:"name"`
}

func FromUser(u *User) UserAsArray {
	return UserAsArray{ID: u.ID, Name: u.Name}
}

// func ToUser(ua UserAsArray) User {
// 	return User{ID: ua.ID, Name: ua.Name}
// }

// func (u User) GetID() uuid.UUID {
// 	return u.ID
// }

func (u *User) SetID(id string) {
	//App.logger.Error("ssdfsdf")
	//todo: log the error
	u.ID, _ = uuid.FromBytes([]byte(id))
}

func (u *User) GenerateID() {
	defer func() { recover() }()

	emptyUUID, _ := uuid.FromBytes(make([]byte, 16))
	if u.ID == emptyUUID {
		u.ID = uuid.New()
	}
}

type Article struct {
	ID     uuid.UUID `json:"id,omitempty"`
	Title  string    `json:"title"`
	Text   string    `json:"text"`
	UserID uuid.UUID `json:"user_id" mapstructure:"user_id"`
}

type ArticleAsArray struct {
	ID     [16]byte `json:"id,omitempty"`
	Title  string   `json:"title"`
	Text   string   `json:"text"`
	UserID [16]byte `json:"user_id"`
}

func FromArticle(a *Article) ArticleAsArray {
	return ArticleAsArray{ID: a.ID, Title: a.Title, Text: a.Text, UserID: a.UserID}
}

// func (a Article) GetID() uuid.UUID {
// 	return a.ID
// }

func (a *Article) SetID(id string) {
	//todo: log the error
	a.ID, _ = uuid.FromBytes([]byte(id))
}

func (a *Article) GenerateID() {
	defer func() { recover() }()

	emptyUUID, _ := uuid.FromBytes(make([]byte, 16))
	if a.ID == emptyUUID {
		a.ID = uuid.New()
	}
}
