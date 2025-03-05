package controllers

import (
	m "DMS/internal/models"
	s "DMS/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHttp struct {
	userService s.UserService
}

func newUserHttp(userService s.UserService) UserHttp {
	return UserHttp{userService}
}

func (h *UserHttp) CreeateUser(c *gin.Context) {
	var user m.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrResponse{
			ErrCode: http.StatusBadRequest,
			Message: BadJsonStruct,
			Details: ParsingError,
		})
		validate.Struct(user)
		fmt.Println(user)
		// id := c.Param("id")
		// c.Status(400)
	}
	h.userService.CreateUser(user.Name, user.PhoneNumber, user.CreatedBy)
}
