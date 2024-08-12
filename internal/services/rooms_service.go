package services

import (
	"context"
	"errors"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/request"
	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	"github.com/JulioZittei/wsrs-ama-go/internal/mappers"
	"github.com/JulioZittei/wsrs-ama-go/internal/repositories"
	"github.com/google/uuid"
)

type RoomsService struct {
	repository    *repositories.RoomsRepository
	roomMapper    *mappers.RoomMapper
	messageMapper *mappers.MessageMapper
}

func NewRoomsService(repository *repositories.RoomsRepository, roomMapper *mappers.RoomMapper,
	messageMapper *mappers.MessageMapper) *RoomsService {
	return &RoomsService{
		repository:    repository,
		roomMapper:    roomMapper,
		messageMapper: messageMapper,
	}
}

func (s *RoomsService) GetRoomMessage(ctx context.Context, roomId uuid.UUID, messageId uuid.UUID) (*response.MessageResponse, error) {
	_, err := s.repository.FindRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}
	message, err := s.repository.FindMessage(ctx, messageId)
	return s.messageMapper.ToResponse(message), err
}

func (s *RoomsService) CreateRoom(ctx context.Context, room *request.RoomRequest) (uuid.UUID, error) {
	return s.repository.SaveRoom(ctx, room.Subject)
}

func (s *RoomsService) GetRooms(ctx context.Context) ([]response.RoomResponse, error) {
	rooms, err := s.repository.FindAllRooms(ctx)
	responseRooms := make([]response.RoomResponse, len(rooms))

	for i, room := range rooms {
		responseRoom := s.roomMapper.ToResponse(&room)
		responseRooms[i] = *responseRoom
	}
	return responseRooms, err
}

func (s *RoomsService) GetRoom(ctx context.Context, roomId uuid.UUID) (*response.RoomResponse, error) {
	room, err := s.repository.FindRoom(ctx, roomId)
	return s.roomMapper.ToResponse(room), err
}

func (s *RoomsService) GetRoomMessages(ctx context.Context, roomId uuid.UUID) ([]response.MessageResponse, error) {
	_, err := s.repository.FindRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}
	messages, err := s.repository.FindAllRoomMessages(ctx, roomId)
	responseMessages := make([]response.MessageResponse, len(messages))

	for i, message := range messages {
		responseMessage := s.messageMapper.ToResponse(&message)
		responseMessages[i] = *responseMessage
	}
	return responseMessages, err
}

func (s *RoomsService) CreateRoomMessage(ctx context.Context, params *request.MessageRequest) (uuid.UUID, error) {
	_, err := s.repository.FindRoom(ctx, params.RoomID)
	if err != nil {
		if errors.Is(err, internal_errors.ErrNotFound) {
			return params.RoomID, internal_errors.NewErrBadRequest(ctx, "room not found")
		}
		return params.RoomID, err
	}
	return s.repository.SaveMessage(ctx, params)
}

func (s *RoomsService) LikeRoomMessage(ctx context.Context, roomId uuid.UUID, messageId uuid.UUID) (int64, error) {
	_, err := s.repository.FindRoom(ctx, roomId)
	if err != nil {
		return 0, err
	}
	return s.repository.LikeMessage(ctx, messageId)
}

func (s *RoomsService) RemoveLikeRoomMessage(ctx context.Context, roomId uuid.UUID, messageId uuid.UUID) (int64, error) {
	_, err := s.repository.FindRoom(ctx, roomId)
	if err != nil {
		return 0, err
	}
	return s.repository.RemoveLikeMessage(ctx, messageId)
}

func (s *RoomsService) AnswerRoomMessage(ctx context.Context, roomId uuid.UUID, messageId uuid.UUID) error {
	_, err := s.repository.FindRoom(ctx, roomId)
	if err != nil {
		return err
	}
	return s.repository.MarkMessageAsAnswered(ctx, messageId)
}
