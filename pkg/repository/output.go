package models

import "github.com/gofiber/fiber/v2"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Output struct {
	Data  interface{} `json:"data"`
	Error []Error     `json:"error"`
}

func NewOutput(data interface{}, errors ...Error) map[string]interface{} {
	return map[string]interface{}{
		"data":  data,
		"error": errors,
	}
}

func NewError(message string, code int) Error {
	return Error{
		Message: message,
		Code:    code,
	}
}
func NewErrors(errors ...error) []Error {
	errs := []Error{}
	for _, err := range errors {
		errs = append(errs, Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}
	return errs
}
