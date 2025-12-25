package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-banking/token"
	"strings"
)

const (
	authorizationHeader     = "Authorization"
	authorizationPrefix     = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(authorizationHeader)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)

		if len(fields) != 2 || fields[0] != authorizationPrefix {
			err := errors.New("authorization header is wrong")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		c.Set(authorizationPayloadKey, payload)

	}

}
