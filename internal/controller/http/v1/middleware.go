package v1

import (
	"net/http"
	"strings"

	"github.com/amiosamu/vk-internship/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	UserIDCtx = "userID"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c.Request)
		if !ok {
			log.Errorf("AuthMiddleware.UserIdentitfy: bearerToken: %v", ErrInvalidAuthHeader)
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			c.Abort()
			return
		}

		userID, err := h.authService.ParseToken(token)
		if err != nil {
			log.Errorf("AuthMiddleware.UserIdentity: bearerToken: %v", ErrInvalidAuthHeader)
			newErrorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader.Error())
			c.Abort()
			return
		}
		c.Set(UserIDCtx, userID)
		c.Next()
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, prefix) {
		return "", false
	}

	return header[len(prefix):], true
}
