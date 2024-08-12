package internal_errors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
	"github.com/JulioZittei/wsrs-ama-go/internal/middlewares"
)

type ErrorBadRequest struct {
	StatusCode int
	StatusText string
	Title      string
	Detail     string
	message    string
}

func NewErrBadRequest(ctx context.Context, detailTag string) *ErrorBadRequest {
	return newErrBadRequest(ctx, detailTag)
}

func newErrBadRequest(ctx context.Context, detailTag string) *ErrorBadRequest {
	statusCode := http.StatusBadRequest
	statusText := strings.ToUpper(http.StatusText(statusCode))
	statusText = strings.Replace(statusText, " ", "_", -1)
	title, _ := locale.GetMessage(ctx, statusText)
	detail, _ := locale.GetMessage(ctx, detailTag)
	message, _ := locale.GetMessage(context.WithValue(ctx, middlewares.LangKey, "en"), detailTag)

	return &ErrorBadRequest{
		StatusCode: statusCode,
		StatusText: statusText,
		Title:      title,
		Detail:     detail,
		message:    message,
	}
}

func (eb *ErrorBadRequest) Error() string {
	return fmt.Sprintf("bad request error: %s", eb.message)
}

var ErrBadRequest = newErrBadRequest(context.Background(), "")
