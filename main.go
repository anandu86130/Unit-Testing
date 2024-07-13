package main

import (
	"userPage/database"
	"userPage/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.CreateDB()
	router := gin.Default()
	routers.UserRoutes(router)
	router.Run(":8081")
}
