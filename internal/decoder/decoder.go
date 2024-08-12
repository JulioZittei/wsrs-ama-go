package decoder

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func DecodeJSON(ctx context.Context, body io.Reader, reqBody interface{}) error {
	if err := render.DecodeJSON(body, reqBody); err != nil {
		return internal_errors.NewErrBadRequest(ctx, "INVALID_JSON")
	}

	validate := validator.New()
	err := validate.Struct(reqBody)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)

	var errorsParam = make([]response.ErrorsParam,
		len(validationErrors))

	for i, v := range validationErrors {
		key := fmt.Sprint(v.Tag())
		field := fmt.Sprint(strings.ToLower(v.Field()[0:1]) + v.Field()[1:])
		paramValue := fmt.Sprint(v.Param())

		message, err := locale.GetMessage(ctx, strings.ToUpper(key), field, paramValue)
		if err != nil {
			return err
		}

		errorsParam[i] = response.ErrorsParam{Param: field, Message: message}
	}

	return internal_errors.NewErrValidation(ctx, errorsParam)

}
