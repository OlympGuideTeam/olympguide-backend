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

	user, complete, err := h.userService.GetUserData(userID)

	if err != nil {
		errs.HandleError(c, err)
		return
	}
	if complete {
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusOK, user.PoorUserDataResponse)
	}
}

// UpdateUser обновляет данные пользователя
// @Summary Обновление профиля
// @Description Обновляет поля профиля: имя, фамилия, отчество, дата рождения, id региона
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param first_name formData string true "Имя"
// @Param last_name formData string true "Фамилия"
// @Param second_name formData string false "Отчество"
// @Param birthday formData string true "Дата рождения (в формате 02.01.2006)"
// @Param region_id formData int true "ID региона"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} errs.AppError "Некорректный формат даты или другие ошибки валидации"
// @Failure 404 {object} errs.AppError "Регион не найден"
// @Failure 500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Security ApiToken
// @Router /user/update [post]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	userID := c.MustGet(constants.ContextUserID).(uint)

	err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
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

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	userID := c.MustGet(constants.ContextUserID).(uint)

	err := h.userService.UpdatePassword(userID, req.Password)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
