package main

import (
	"apollo-proxy/config"
	"apollo-proxy/controller"
	"apollo-proxy/middleware"
	"apollo-proxy/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.Config.App.Mode)
	router := gin.Default()
	router.Use(middleware.Cors(), middleware.Auth())
	apollo := new(controller.Apollo)
	router.GET("/:env/:project/:branch/configs/:apollo_app_id/:cluster_name/:namespace", apollo.Configs)
	router.GET("/:env/:project/:branch/notifications/v2", apollo.Notifications)
	go func() {
		service.Watch()
	}()
	router.Run(fmt.Sprintf(":%d", config.Config.App.Port))
}
