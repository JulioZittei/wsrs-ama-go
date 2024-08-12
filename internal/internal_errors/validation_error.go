package internal_errors

import (
	"context"
	"net/http"
	"strings"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
)

type ErrorValidation struct {
	StatusCode  int
	StatusText  string
	Title       string
	Detail      string
	ErrorsParam []response.ErrorsParam
}

func NewErrValidation(ctx context.Context, errorsParam []response.ErrorsParam) *ErrorValidation {
	statusCode := http.StatusUnprocessableEntity
	statusText := strings.ToUpper(http.StatusText(statusCode))
	statusText = strings.Replace(statusText, " ", "_", -1)
	title, _ := locale.GetMessage(ctx, statusText)

	return &ErrorValidation{
		StatusCode:  statusCode,
		StatusText:  statusText,
		Title:       title,
		ErrorsParam: errorsParam,
	}
}

func (ev *ErrorValidation) Error() string {
	return "validation error"
}

var ErrValidation = NewErrValidation(context.Background(), nil)
