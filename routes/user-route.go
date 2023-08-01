package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine){
	//All routes related to user comes here
	router.POST("/user", controllers.CreatePet())
	router.GET("/user/:userId", controllers.GetPet())
	router.GET("/users", controllers.GetAllPets())
	router.PUT("/user/:userId", controllers.EditAPet())
	router.DELETE("/user/:userId", controllers.DeleteAPet())
}