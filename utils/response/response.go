package response

import "github.com/go-playground/validator/v10"

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string            `json:"code,omitempty" example:"VALIDATION_ERROR"`
	Message string            `json:"message" example:"Invalid request body"`
	Errors  map[string]string `json:"errors,omitempty" swaggertype:"object"`
}

// SuccessResponse represents a success response with optional data
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
	Data    any    `json:"data,omitempty" swaggertype:"object"`
}

// Success sends a success response with a message and optional data
func Success(message string, data any) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}

}

// SuccessMessage sends a success response with only a message
func SuccessMessage(message string) SuccessResponse {
	return SuccessResponse{
		Message: message,
	}
}

// ValidationErrorResponse builds a structured error response from validator errors.
func ValidationErrorResponse(err error, message string, code string) (ErrorResponse, bool) {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse{}, false
	}

	fieldErrors := make(map[string]string, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := fieldError.Field()
		if field == "" {
			field = fieldError.StructField()
		}

		tag := fieldError.Tag()
		if tag == "" {
			tag = fieldError.ActualTag()
		}

		if field != "" {
			fieldErrors[field] = tag
		}
	}

	return ErrorResponse{
		Code:    code,
		Message: message,
		Errors:  fieldErrors,
	}, true
}
