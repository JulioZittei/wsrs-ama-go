package repositories

import (
	"context"
	"log/slog"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/request"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	"github.com/JulioZittei/wsrs-ama-go/internal/mappers"
	"github.com/JulioZittei/wsrs-ama-go/internal/models"
	"github.com/JulioZittei/wsrs-ama-go/internal/store/pgstore"
	"github.com/google/uuid"
)

type RoomsRepository struct {
	db            *pgstore.Queries
	roomMapper    *mappers.RoomMapper
	messageMapper *mappers.MessageMapper
}

func NewRoomsRepository(db *pgstore.Queries, roomMapper *mappers.RoomMapper,
	messageMapper *mappers.MessageMapper) *RoomsRepository {
	return &RoomsRepository{
		db:            db,
		roomMapper:    roomMapper,
		messageMapper: messageMapper,
	}
}

func (rr *RoomsRepository) FindMessage(ctx context.Context, messageId uuid.UUID) (*models.Message, error) {
	message, err := rr.db.GetMessage(ctx, messageId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			slog.Error("message not found", "error", err)
			return rr.messageMapper.ToModel(message), internal_errors.NewErrNotFound(ctx, "Message")
		}

		slog.Error("something went wrong while finding a message", "error", err)
		return rr.messageMapper.ToModel(message), internal_errors.NewErrInternal(ctx, err)
	}
	return rr.messageMapper.ToModel(message), err
}

func (rr *RoomsRepository) FindRoom(ctx context.Context, roomId uuid.UUID) (*models.Room, error) {
	room, err := rr.db.GetRoom(ctx, roomId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			slog.Error("room not found", "error", err)
			return rr.roomMapper.ToModel(room), internal_errors.NewErrNotFound(ctx, "Room")
		}

		slog.Error("something went wrong while finding a room", "error", err)
		return rr.roomMapper.ToModel(room), internal_errors.NewErrInternal(ctx, err)
	}
	return rr.roomMapper.ToModel(room), err
}

func (rr *RoomsRepository) FindAllRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := rr.db.GetRooms(ctx)
	modelRooms := make([]models.Room, len(rooms))

	for i, room := range rooms {
		modelRoom := rr.roomMapper.ToModel(room)
		modelRooms[i] = *modelRoom
	}

	if err != nil {
		slog.Error("something went wrong while finding all rooms", "error", err)
		return modelRooms, internal_errors.NewErrInternal(ctx, err)
	}

	return modelRooms, err
}

func (rr *RoomsRepository) FindAllRoomMessages(ctx context.Context, roomID uuid.UUID) ([]models.Message, error) {
	messages, err := rr.db.GetRoomMessages(ctx, roomID)
	modelMessages := make([]models.Message, len(messages))

	for i, message := range messages {
		modelMessage := rr.messageMapper.ToModel(message)
		modelMessages[i] = *modelMessage
	}

	if err != nil {
		slog.Error("something went wrong while finding all room messages", "error", err)
		return modelMessages, internal_errors.NewErrInternal(ctx, err)
	}
	return modelMessages, err
}

func (rr *RoomsRepository) SaveRoom(ctx context.Context, subject string) (uuid.UUID, error) {
	roomId, err := rr.db.InsertRoom(ctx, subject)
	if err != nil {
		slog.Error("something went wrong while saving room", "error", err)
		return roomId, internal_errors.NewErrInternal(ctx, err)
	}
	return roomId, err
}

func (rr *RoomsRepository) SaveMessage(ctx context.Context, params *request.MessageRequest) (uuid.UUID, error) {
	messageId, err := rr.db.InsertMessage(ctx, pgstore.InsertMessageParams{
		RoomID:  params.RoomID,
		Message: params.Message,
	})
	if err != nil {
		slog.Error("something went wrong while saving message", "error", err)
		return messageId, internal_errors.NewErrInternal(ctx, err)
	}
	return messageId, err
}

func (rr *RoomsRepository) LikeMessage(ctx context.Context, messageId uuid.UUID) (int64, error) {
	likeCount, err := rr.db.ReactToMessage(ctx, messageId)
	if err != nil {
		slog.Error("something went wrong while adding like message", "error", err)
		return likeCount, internal_errors.NewErrInternal(ctx, err)
	}
	return likeCount, err
}

func (rr *RoomsRepository) RemoveLikeMessage(ctx context.Context, messageId uuid.UUID) (int64, error) {
	likeCount, err := rr.db.RemoveReactionFromMessage(ctx, messageId)
	if err != nil {
		slog.Error("something went wrong while removing like message", "error", err)
		return likeCount, internal_errors.NewErrInternal(ctx, err)
	}
	return likeCount, err
}

func (rr *RoomsRepository) MarkMessageAsAnswered(ctx context.Context, messageId uuid.UUID) error {
	err := rr.db.MarkMessageAsAnswered(ctx, messageId)
	if err != nil {
		slog.Error("something went wrong while marking message as answered", "error", err)
		return internal_errors.NewErrInternal(ctx, err)
	}
	return nil
}
