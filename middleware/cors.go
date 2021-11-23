package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Cors() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			log.Println("OPTIONS")
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			ctx.AbortWithStatus(200)
			return
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Next()
	}
}
