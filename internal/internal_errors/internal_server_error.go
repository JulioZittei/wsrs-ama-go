package internal_errors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
)

type ErrorInternalServer struct {
	StatusCode int
	StatusText string
	Title      string
	Detail     string
	cause      error
}

func (ei *ErrorInternalServer) Error() string {
	return fmt.Sprintf("internal server error: %s", ei.cause.Error())
}

func NewErrInternal(ctx context.Context, err error) *ErrorInternalServer {
	statusCode := http.StatusInternalServerError
	statusText := strings.ToUpper(http.StatusText(statusCode))
	statusText = strings.Replace(statusText, " ", "_", -1)
	title, _ := locale.GetMessage(ctx, statusText)

	return &ErrorInternalServer{
		StatusCode: statusCode,
		StatusText: statusText,
		Title:      title,
		cause:      err,
	}
}

var ErrInternalServer = NewErrInternal(context.Background(), nil)
