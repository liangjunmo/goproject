package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjunmo/gocode"

	"github.com/liangjunmo/goproject/internal/app/server/config"
	"github.com/liangjunmo/goproject/internal/app/server/servercode"
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

	code := gocode.Parse(err)
	if code == gocode.SuccessCode {
		code = servercode.OK
	} else if code == gocode.DefaultCode {
		code = servercode.Unknown
	}

	body := gin.H{
		"data": data,
		"code": code,
		"msg":  servercode.Translate(code, servercode.Language(c.GetHeader("Accept-Language"))),
	}

	if config.Config.Debug {
		body["error"] = nil
		if err != nil {
			body["error"] = err.Error()
		}

		body["request_id"] = c.Request.Context().Value(config.TraceIdKey)
	}

	return body
}

func (handler *BaseHandler) GetUserClaims(c *gin.Context) *UserJwtClaims {
	user, _ := c.Get("user_claims")
	return user.(*UserJwtClaims)
}
