package mappers

import (
	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/models"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
)

type MessageMapper struct{}

func (mapper *MessageMapper) ToModel(message pgstore.Message) *models.Message {
	return &models.Message{
		ID:         message.ID,
		RoomID:     message.RoomID,
		Message:    message.Message,
		LikesCount: message.LikesCount,
		Answered:   message.Answered,
	}
}

func (mapper *MessageMapper) ToResponse(message *models.Message) *response.MessageResponse {
	return &response.MessageResponse{
		ID:         message.ID.String(),
		RoomID:     message.RoomID.String(),
		Message:    message.Message,
		LikesCount: message.LikesCount,
		Answered:   message.Answered,
	}
}
