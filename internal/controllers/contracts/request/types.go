package request

import "github.com/google/uuid"

type RoomRequest struct {
	Subject string `json:"subject" validate:"required"`
}

type MessageRequest struct {
	RoomID  uuid.UUID
	Message string `json:"message" validate:"required"`
}
