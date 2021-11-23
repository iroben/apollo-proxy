package controller

import (
	"github.com/gin-gonic/gin"
)

var (
	RESP_OK        = 10000
	RESP_ERR       = 10001
	RESP_PARA_ERR  = 10002
	RESP_NOT_LOGIN = 10003
)
var msg = map[int]string{
	RESP_OK:        "成功",
	RESP_ERR:       "请求失败",
	RESP_PARA_ERR:  "参数验证失败",
	RESP_NOT_LOGIN: "未登录",
}

const (
	TEST_USER_ID   = "liubin10"
	TEST_USER_NAME = "刘斌"
)

type BaseController struct {
}
type UserInfo struct {
	UserId   string
	UserName string
}

func (c *BaseController) GetUserInfo() *UserInfo {
	return &UserInfo{
		UserId:   TEST_USER_ID,
		UserName: TEST_USER_NAME,
	}
}
func (c *BaseController) JSON(ctx *gin.Context, data map[string]interface{}) {
	if _, ok := data["msg"]; ok {
		ctx.JSON(200, data)
		return
	}
	if code, ok := data["code"]; ok {
		if codeInt, ok := code.(int); ok {
			data["msg"] = msg[codeInt]
		}
	}
	ctx.JSON(200, data)
}
