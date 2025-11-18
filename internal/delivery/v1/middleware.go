package v1

import (
	"avito-internship/pkg/e"
	"avito-internship/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	logger     logger.Logger
	adminToken string
}

func NewMiddleware(logger logger.Logger, adminToken string) *Middleware {
	return &Middleware{
		logger:     logger,
		adminToken: adminToken,
	}
}

//func (m *Middleware) AdminMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		const prefix = "Bearer "
//		authHeader := c.GetHeader("Authorization")
//		token := strings.TrimPrefix(authHeader, prefix)
//
//		if !strings.HasPrefix(authHeader, prefix) || token != m.adminToken {
//			m.logger.Errorf(e.ErrUnauthorized, "method=%s path=%s", c.Request.Method, c.Request.URL.Path)
//			abortUnauthorized(c)
//			return
//		}
//
//		c.Next()
//	}
//}

func (m *Middleware) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			err := ginErr.Err

			m.logger.Errorf(err, "method=%s path=%s", c.Request.Method, c.Request.URL.Path)

			codeInt, codeString, msg := ToHTTPResponse(err)
			response := NewErrorResponse(codeString, msg)
			c.JSON(codeInt, response)
			return
		}
	}
}

func abortUnauthorized(c *gin.Context) {
	response := NewErrorResponse(e.NOT_FOUND, e.ErrResourceNotFound.Error())
	c.JSON(http.StatusUnauthorized, response)
	c.Abort()
}
