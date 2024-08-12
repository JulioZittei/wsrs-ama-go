package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	locale "github.com/JulioZittei/wsrs-ama-go/internal/locale/message_bundle"
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(ctx context.Context, obj interface{}) error {
	validate := validator.New()
	err := validate.Struct(obj)
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
