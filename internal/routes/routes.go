package routes

import (
	c "DMS/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// In general, it's better to return json as the details field of the response if
// the response http code is 200.
func SetupRouter(router *gin.Engine, ctr c.HttpConrtoller) {
	router.Use(ctr.Middleware.Cors)

	apiV1NeedAuth(router, ctr)
	apiV1NeedNotAuth(router, ctr)

	router.GET("/health", healthCheck)
	// Open this path to see documentaion=> /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
func apiV1NeedAuth(router *gin.Engine, ctr c.HttpConrtoller) {
	routerV1 := router.Group("api/v1")
	routerV1.Use(ctr.Middleware.Authentication)

	routerV1.POST("/users", ctr.User.CreateUser)
	routerV1.GET("/users/current", ctr.User.GetCurrentUserInfo)
	routerV1.POST("/jps", ctr.JP.CreateUserJP)
	routerV1.POST("/jps/admin", ctr.JP.CreateAdminJP)
	routerV1.GET("/user/jps", ctr.JP.GetUserJPs)
	// Create an event.
	// If response http code be 200, then return json as details field of the response.
	routerV1.POST("/events", ctr.Event.CreateEvent)
	routerV1.GET("/events", ctr.Event.GetNLastEventsByJPID)
	routerV1.POST("/docs", ctr.Doc.CreateDoc)
	routerV1.GET("/docs", ctr.Doc.GetNLastDocs)
	routerV1.GET("/jps/:jp_id/events/:event_id/docs", ctr.Doc.GetNLastDocsByEventID)
	routerV1.POST("/logout", ctr.Session.Logout)
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
