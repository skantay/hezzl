package ginV1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skantay/hezzl/internal/domain/good/model"
)

func badRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func internalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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

func handleError(c *gin.Context, err error) {
	if errors.Is(err, model.ErrGoodNotFound) {
		type response struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Details string `json:"details"`
		}

		resp := response{
			Code:    3,
			Message: err.Error(),
		}

		c.JSON(http.StatusNotFound, resp)
	} else {
		internalServerError(c)
	}
}
