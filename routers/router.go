package routers

import (
	"userPage/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/signup", controller.Signup)
	router.POST("/signin", controller.SignIn)
	router.GET("/userlist", controller.UserList)
	router.PATCH("/user/edit/:id", controller.EditUser)

	router.POST("admin/login", controller.AdminLogin)
}
