package v1

import (
	"github.com/gin-gonic/gin"
)

type DefaultHandler struct {
	*BaseHandler
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (handler *DefaultHandler) Health(c *gin.Context) {
	c.String(200, "healthy")
}
