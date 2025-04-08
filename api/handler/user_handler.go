package handler

import (
	"api/dto"
	"api/service"
	"api/utils/constants"
	"api/utils/errs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserData(c *gin.Context) {
	userID, _ := c.MustGet(constants.ContextUserID).(uint)

	user, err := h.userService.GetUserData(userID)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uint)
	err := h.userService.DeleteUser(userID)
	if err != nil {
		errs.HandleError(c, err)
	}

	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: constants.AccountDeleted})
}
