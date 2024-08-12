package internal_errors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
)

type ErrorNotFound struct {
	StatusCode int
	StatusText string
	Title      string
	Detail     string
	resource   string
	message    string
}

func NewErrNotFound(ctx context.Context, resource string) *ErrorNotFound {
	if resource != "" {
		return newErrNotFound(ctx, resource)
	}
	return newErrNotFound(ctx, "resource")
}

func newErrNotFound(ctx context.Context, resource string) *ErrorNotFound {
	statusCode := http.StatusNotFound
	statusText := strings.ToUpper(http.StatusText(statusCode))
	statusText = strings.Replace(statusText, " ", "_", -1)
	title, _ := locale.GetMessage(ctx, statusText)
	detail, _ := locale.GetMessage(ctx, "RESOURCE_NOT_FOUND", strings.ToLower(resource))
	message, _ := locale.GetMessage(context.Background(), "RESOURCE_NOT_FOUND", resource)

	return &ErrorNotFound{
		StatusCode: statusCode,
		StatusText: statusText,
		Title:      title,
		Detail:     detail,
		resource:   resource,
		message:    message,
	}
}

func (en *ErrorNotFound) Error() string {
	return fmt.Sprintf("resource not found error: %s", en.message)
}

var ErrNotFound = newErrNotFound(context.Background(), "")
