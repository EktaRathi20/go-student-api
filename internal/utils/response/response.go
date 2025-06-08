package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errorMessages []string
	for _, fieldError := range errors {
		switch fieldError.ActualTag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is required", fieldError.Field()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid", fieldError.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  "Validation failed: " + strings.Join(errorMessages, ", "),
	}
}
