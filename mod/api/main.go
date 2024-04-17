package main

import (
	"github.com/berachain/beacon-kit/mod/api/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())

	router.POST("/genesis", handler.AddAccountAndPredeploy)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
