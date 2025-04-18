package routes

import (
	"github.com/gin-gonic/gin"
	"bigredlink/controllers"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)
	r.POST("/logout", controllers.Logout)
}
