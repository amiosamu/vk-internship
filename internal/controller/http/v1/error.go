package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
)




var (
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrCannotParseToken = errors.New("cannot parse token")
)


type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	err := errors.New(message)
	_, ok := err.(*gin.Error)

	if !ok {
		report := gin.H {
			"error" : err.Error(),
		}
		c.JSON(statusCode, report)
	}
	c.Error(errors.New("internal server error"))
}