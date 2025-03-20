package routes

import (
	c "DMS/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// In general, it's better to return json as the details field of the response if
// the response http code is 200.
func SetupRouter(router *gin.Engine, ctr c.HttpConrtoller) {
	apiV1NeedAuth(router, ctr)
	apiV1NeedNotAuth(router, ctr)

	router.GET("/health", healthCheck)
}

func apiV1NeedAuth(router *gin.Engine, ctr c.HttpConrtoller) {
	routerV1 := router.Group("api/v1")
	routerV1.Use(ctr.Middleware.Authentication)

	routerV1.POST("/users", ctr.User.CreateUser)
	routerV1.POST("/jps", ctr.JP.CreateJP)
	// Return job position list of the specified user id
	routerV1.GET("/users/:id/jp", ctr.JP.GetUserJPs)
	// Create an event.
	// If response http code be 200, then return json as details field of the response.
	routerV1.POST("/events", ctr.Event.CreateEvent)
	routerV1.POST("/docs", ctr.Doc.CreateDoc)
	// This route have the "count" query parameter
	routerV1.GET("/docs/event/:id", ctr.Doc.GetNLastDocsByEventID)
	// router.GET("/users/:id", controller.GetUser)
	// router.GET("/products", controllers.GetProducts) //Example of a different controller.
}

func apiV1NeedNotAuth(router *gin.Engine, ctr c.HttpConrtoller) {
	routerV1 := router.Group("api/v1")

	// Admin is a sub-type of user entity
	routerV1.POST("/users/admin", ctr.User.CreateAdmin)
	routerV1.POST("login/phone-based", ctr.Session.PhoneBasedLogin)
}

// Check the healthy status of services
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
