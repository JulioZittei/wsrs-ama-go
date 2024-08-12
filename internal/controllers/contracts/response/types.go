package response

type RoomResponse struct {
	ID      string `json:"id"`
	Subject string `json:"subject,omitempty"`
}

type MessageResponse struct {
	ID         string `json:"id"`
	RoomID     string `json:"room_id"`
	Message    string `json:"message,omitempty"`
	LikesCount int64  `json:"likes_count,omitempty"`
	Answered   bool   `json:"is_answered,omitempty"`
}

type ErrorsParam struct {
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Status        string        `json:"status,omitempty"`
	Code          int           `json:"code,omitempty"`
	Title         string        `json:"title,omitempty"`
	Detail        string        `json:"detail,omitempty"`
	Instance      string        `json:"instance,omitempty"`
	InvalidParams []ErrorsParam `json:"invalid_params,omitempty"`
}
