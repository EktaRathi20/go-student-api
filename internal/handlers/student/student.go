package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"student-api/internal/types"
	"student-api/internal/utils/response"

	"github.com/go-playground/validator"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		error := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(error, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if error != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(error))
			return
		}

		// Request validation
		if err := validator.New().Struct(&student); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrors))
			return
		}
		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
