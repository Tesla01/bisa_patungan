package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

func GetEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("No ENV File")
	}

	return os.Getenv(key)
}
