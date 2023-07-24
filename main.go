package main

import "github.com/gin-gonic/gin"


func main() {
  router := gin.Default()
	
  router.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "data": "Hello from Gin-gonic & mongoDB",
    })
  })
  router.Run("localhost:6000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}