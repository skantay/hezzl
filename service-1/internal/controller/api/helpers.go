package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type responseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func parseQueryParamAtoi(c *gin.Context, paramName string, defaultValue int) (int, error) {
	value := c.Query(paramName)
	if value == "" {
		return defaultValue, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func handleError(c *gin.Context, details string, code int, err error) {
	var msg string

	if err != nil {
		msg = err.Error()
	}

	resp := responseError{
		Code:    code,
		Message: msg,
		Details: details,
	}

	c.JSON(code, resp)
}
