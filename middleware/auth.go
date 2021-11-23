package middleware

import (
	"apollo-proxy/config"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userInfo := ""
		if _, ok := ctx.Request.Header["Authorization"]; !ok {
			ctx.AbortWithStatus(401)
			return
		}
		basicStr := ctx.Request.Header["Authorization"][0]
		authorizationData := strings.Split(basicStr, " ")
		_userInfo, _ := base64.StdEncoding.DecodeString(authorizationData[1])
		userInfo = string(_userInfo)
		userInfoArray := strings.Split(userInfo, ":")

		if userInfoArray[0] == config.Config.App.User &&
			userInfoArray[1] == config.Config.App.Passwd {
			ctx.Next()
			return
		}
		ctx.AbortWithStatus(401)
	}
}
