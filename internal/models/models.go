package models

import "github.com/google/uuid"

type Message struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	Message    string
	LikesCount int64
	Answered   bool
}

type Room struct {
	ID      uuid.UUID
	Subject string
}
