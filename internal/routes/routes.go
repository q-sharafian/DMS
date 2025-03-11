package routes

import (
	c "DMS/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, ctr c.HttpConrtoller) {
	router.POST("/users", ctr.User.CreateUser)
	// Admin is a sub-type of user entity
	router.POST("/users/admin", ctr.User.CreateAdmin)
	// router.POST("/jps", ctr.JP.CreateJP)

	// router.GET("/users/:id", controller.GetUser)
	// router.GET("/products", controllers.GetProducts) //Example of a different controller.

	router.GET("/health", healthCheck)
}

// Check the healthy status of services
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
