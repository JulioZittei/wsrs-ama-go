package mappers

import (
	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/models"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
)

type RoomMapper struct{}

func (mapper *RoomMapper) ToModel(room pgstore.Room) *models.Room {
	return &models.Room{
		ID:      room.ID,
		Subject: room.Subject,
	}
}

func (mapper *RoomMapper) ToResponse(room *models.Room) *response.RoomResponse {
	return &response.RoomResponse{
		ID:      room.ID.String(),
		Subject: room.Subject,
	}
}
