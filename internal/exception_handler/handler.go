package exception_handler

import (
	"net/http"

	"github.com/JulioZittei/wsrs-ama-go/internal/controllers/contracts/response"
	"github.com/JulioZittei/wsrs-ama-go/internal/internal_errors"
	"github.com/go-chi/render"
)

type ControllerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func ExceptionHandler(controllerFunc ControllerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj, status, err := controllerFunc(w, r)

		if err != nil {
			errorResponse := handleError(r, err)
			w.WriteHeader(status)
			render.JSON(w, r, errorResponse)
			return
		}

		w.WriteHeader(status)
		if obj != nil {
			render.JSON(w, r, obj)
		}
	})
}

func handleError(r *http.Request, err error) *response.ErrorResponse {
	switch err := err.(type) {
	case *internal_errors.ErrorValidation:
		return buildErrValidationResponse(r, err)
	case *internal_errors.ErrorInternalServer:
		return buildErrInternalResponse(r, err)
	case *internal_errors.ErrorBadRequest:
		return buildBadRequestResponse(r, err)
	case *internal_errors.ErrorNotFound:
		return buildNotFoundResponse(r, err)
	default:
		return buildDefaultErrorResponse(r)
	}

}

func buildErrValidationResponse(r *http.Request, err *internal_errors.ErrorValidation) *response.ErrorResponse {
	return &response.ErrorResponse{
		Code:          err.StatusCode,
		Status:        http.StatusText(err.StatusCode),
		Title:         err.Title,
		Detail:        err.Detail,
		Instance:      r.RequestURI,
		InvalidParams: err.ErrorsParam,
	}
}

func buildErrInternalResponse(r *http.Request, err *internal_errors.ErrorInternalServer) *response.ErrorResponse {
	return &response.ErrorResponse{
		Code:          err.StatusCode,
		Status:        http.StatusText(err.StatusCode),
		Title:         err.Title,
		Detail:        err.Detail,
		Instance:      r.RequestURI,
		InvalidParams: []response.ErrorsParam{},
	}
}

func buildBadRequestResponse(r *http.Request, err *internal_errors.ErrorBadRequest) *response.ErrorResponse {
	return &response.ErrorResponse{
		Code:          err.StatusCode,
		Status:        http.StatusText(err.StatusCode),
		Title:         err.Title,
		Detail:        err.Detail,
		Instance:      r.RequestURI,
		InvalidParams: []response.ErrorsParam{},
	}
}

func buildNotFoundResponse(r *http.Request, err *internal_errors.ErrorNotFound) *response.ErrorResponse {
	return &response.ErrorResponse{
		Code:          err.StatusCode,
		Status:        http.StatusText(err.StatusCode),
		Title:         err.Title,
		Detail:        err.Detail,
		Instance:      r.RequestURI,
		InvalidParams: []response.ErrorsParam{},
	}
}

func buildDefaultErrorResponse(r *http.Request) *response.ErrorResponse {
	return &response.ErrorResponse{
		Code:          http.StatusInternalServerError,
		Status:        http.StatusText(http.StatusInternalServerError),
		Title:         "unmaped error",
		Detail:        "contact support to report the problem",
		Instance:      r.RequestURI,
		InvalidParams: []response.ErrorsParam{},
	}
}
