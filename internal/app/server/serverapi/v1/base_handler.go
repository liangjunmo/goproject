package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gocode"

	"github.com/liangjunmo/goproject/internal/server/servercode"
)

type BaseHandler struct{}

func (handler *BaseHandler) Response(c *gin.Context, data interface{}, err error) {
	c.JSON(200, handler.buildResponseBody(c, data, err))
}

func (handler *BaseHandler) ResponseWithStatusCode(c *gin.Context, statusCode int, data interface{}, err error) {
	c.JSON(statusCode, handler.buildResponseBody(c, data, err))
}

func (handler *BaseHandler) buildResponseBody(c *gin.Context, data interface{}, err error) gin.H {
	if data == nil {
		data = map[string]interface{}{}
	}

	lang := c.GetHeader("Accept-Language")
	if _, ok := i18n[lang]; !ok {
		lang = "zh_CN"
	}

	code := gocode.Parse(err)
	if code == gocode.SuccessCode {
		code = servercode.OK
	} else if code == gocode.DefaultCode {
		code = servercode.Unknown
	}

	body := gin.H{
		"data": data,
		"code": code,
		"msg":  i18n[lang][code],
	}

	return body
}
