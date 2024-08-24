package request

import "github.com/google/uuid"

type RoomRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type MessageRequest struct {
	RoomID  uuid.UUID `json:"-"`
	Message string    `json:"message" validate:"required"`
}
