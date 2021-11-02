package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"loginWithEmailGo/login"
)

func main()  {
	router := gin.Default()
	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Allow-Methods", "*")
		if context.Request.Method == "OPTIONS" {
			context.Status(200)
			context.Abort()
		}
	})

	loginGroup := router.Group("/user")
	login.LoginSetup(loginGroup)

	err := router.Run(":8000")
	if err != nil {
		fmt.Println("Connection can not be completed!")
		return
	}
}
