package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine){
	//All routes related to user comes here
	router.POST("/user", controllers.CreateUser())
}