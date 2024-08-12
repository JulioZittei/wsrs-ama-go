package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/request"
	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/socket"
	"github.com/JulioZittei/wsrs-ama-go/internal/decoder"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	"github.com/JulioZittei/wsrs-ama-go/internal/services"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomsController struct {
	service     *services.RoomsService
	upgrader    websocket.Upgrader
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	mutex       *sync.Mutex
}

func NewRoomsController(service *services.RoomsService, upgrader websocket.Upgrader) *RoomsController {
	return &RoomsController{
		service:     service,
		upgrader:    upgrader,
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mutex:       &sync.Mutex{},
	}
}

func (c *RoomsController) CreateRoom(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var requestBody = request.RoomRequest{}
	if err := decoder.DecodeJSON(r.Context(), r.Body, &requestBody); err != nil {
		return nil, 400, err
	}

	roomId, err := c.service.CreateRoom(r.Context(), &requestBody)
	if err != nil {
		return nil, 500, err
	}
	data := &response.RoomResponse{
		ID:      roomId.String(),
		Subject: requestBody.Subject,
	}

	return data, 201, err
}

func (c *RoomsController) CreateRoomMessage(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return nil, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	var requestBody = request.MessageRequest{}
	if err := decoder.DecodeJSON(r.Context(), r.Body, &requestBody); err != nil {
		return nil, 400, err
	}

	requestBody.RoomID = roomId

	messageId, err := c.service.CreateRoomMessage(r.Context(), &requestBody)
	if err != nil {
		return nil, 500, err
	}

	data := &response.MessageResponse{
		ID:      messageId.String(),
		RoomID:  rawRoomId,
		Message: requestBody.Message,
	}

	go c.notifyClients(socket.Message{
		Kind: socket.MessageKindMessageCreated,

		RoomID: rawRoomId,
		Value: socket.MessageMessageCreated{
			ID:      data.ID,
			Message: data.Message,
		},
	})

	return data, 201, err
}

func (c *RoomsController) GetRooms(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	room, err := c.service.GetRooms(r.Context())
	return room, 200, err
}

func (c *RoomsController) GetRoomMessages(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return nil, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	rooms, err := c.service.GetRoomMessages(r.Context(), roomId)
	return rooms, 200, err
}

func (c *RoomsController) GetRoom(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return nil, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}
	room, err := c.service.GetRoom(r.Context(), roomId)
	return room, 200, err
}

func (c *RoomsController) GetRoomMessage(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return nil, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	rawMessageId := chi.URLParam(r, "message_id")
	messageId, err := uuid.Parse(rawMessageId)
	if err != nil {
		return nil, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid message id")
	}

	message, err := c.service.GetRoomMessage(r.Context(), roomId, messageId)

	return message, 200, err
}

func (c *RoomsController) LikeRoomMessage(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	rawMessageId := chi.URLParam(r, "message_id")
	messageId, err := uuid.Parse(rawMessageId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid message id")
	}
	likeCount, err := c.service.LikeRoomMessage(r.Context(), roomId, messageId)

	go c.notifyClients(socket.Message{
		Kind:   socket.MessageKindMessageRactionIncreased,
		RoomID: rawRoomId,
		Value: socket.MessageMessageReactionIncreased{
			ID:    rawMessageId,
			Count: likeCount,
		},
	})

	return likeCount, 200, err
}
func (c *RoomsController) RemoveLikeRoomMessage(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	rawMessageId := chi.URLParam(r, "message_id")
	messageId, err := uuid.Parse(rawMessageId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid message id")
	}
	likeCount, err := c.service.RemoveLikeRoomMessage(r.Context(), roomId, messageId)

	go c.notifyClients(socket.Message{
		Kind:   socket.MessageKindMessageRactionDecreased,
		RoomID: rawRoomId,
		Value: socket.MessageMessageReactionDecreased{
			ID:    rawMessageId,
			Count: likeCount,
		},
	})

	return likeCount, 200, err
}

func (c *RoomsController) AnswerRoomMessage(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid room id")
	}

	rawMessageId := chi.URLParam(r, "message_id")
	messageId, err := uuid.Parse(rawMessageId)
	if err != nil {
		return 0, 400, internal_errors.NewErrBadRequest(r.Context(), "invalid message id")
	}

	go c.notifyClients(socket.Message{
		Kind:   socket.MessageKindMessageAnswered,
		RoomID: rawRoomId,
		Value: socket.MessageMessageAnswered{
			ID: rawMessageId,
		},
	})

	return nil, 200, c.service.AnswerRoomMessage(r.Context(), roomId, messageId)
}

func (c *RoomsController) notifyClients(message socket.Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	subscribers, ok := c.subscribers[message.RoomID]
	if !ok || len(subscribers) == 0 {
		return
	}

	for conn, cancel := range subscribers {
		if err := conn.WriteJSON(message); err != nil {
			slog.Error("failed to send message to client", "error", err)
			cancel()
		}
	}
}

func (c *RoomsController) SubscribeRoom(w http.ResponseWriter, r *http.Request) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	_, err = c.service.GetRoom(r.Context(), roomId)
	if err != nil {
		if errors.Is(err, internal_errors.ErrNotFound) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("failed to upgrade connection", "error", err)
		http.Error(w, "failed to upgrade websocket connection", http.StatusBadRequest)
		return
	}

	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())

	c.mutex.Lock()
	if _, ok := c.subscribers[rawRoomId]; !ok {
		c.subscribers[rawRoomId] = make(map[*websocket.Conn]context.CancelFunc)
	}
	slog.Info("new client connected", "room_id", rawRoomId, "client_ip", r.RemoteAddr)
	c.subscribers[rawRoomId][conn] = cancel
	c.mutex.Unlock()

	<-ctx.Done()

	c.mutex.Lock()
	delete(c.subscribers[rawRoomId], conn)
	c.mutex.Unlock()
}
