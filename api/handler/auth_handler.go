package handler

import (
	"api/dto"
	"api/service"
	"api/utils/constants"
	"api/utils/errs"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService       service.IAuthService
	googleAuthService service.IGoogleAuthService
	tokenService      service.ITokenService
}

func NewAuthHandler(
	authService service.IAuthService,
	googleAuthService service.IGoogleAuthService,
	tokenService service.ITokenService,
) *AuthHandler {
	return &AuthHandler{authService: authService, googleAuthService: googleAuthService, tokenService: tokenService}
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var request dto.SendRequest
	if err := c.ShouldBind(&request); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	err := h.authService.SendCode(request.Email)

	if err != nil {
		errs.HandleError(c, err)
		return
	}

	log.Printf("Code sent to %s", request.Email)
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Code sent to " + request.Email})
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var request dto.VerifyRequest
	if err := c.ShouldBind(&request); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	err := h.authService.VerifyCode(request.Email, request.Code)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	tempToken, err := h.tokenService.GenerateEmailToken(request.Email)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.VerifyCodeResponse{
		Message: constants.EmailConfirmed,
		Token:   tempToken,
	})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var request dto.SignUpRequest
	if err := c.ShouldBind(&request); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}
	request.Email = c.MustGet(constants.ContextEmail).(string)
	fmt.Println("Email: ", request.Email)

	err := h.authService.SignUp(&request)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: constants.SignedUp})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	user, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set(constants.ContextUserID, user.UserID)

	if err = session.Save(); err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message:   constants.LoggedIn,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		errs.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: constants.LoggedOut})
}

func (h *AuthHandler) CheckSession(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get(constants.ContextUserID)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{Message: constants.Unauthorized})
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: constants.Authorized})
}

// GoogleLogin выполняет вход через Google OAuth2.
//
// @Summary Вход через Google
// @Description Если пользователь уже завершил регистрацию, создаётся сессия. Иначе возвращается временный токен для завершения регистрации.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.GoogleAuthRequest true "Токен Google ID"
// @Success 200 {object} dto.LoginResponse "Регистрация завершена — пользователь вошёл"
// @Success 200 {object} dto.RegistrationIncompleteResponse "Регистрация не завершена — вернётся временный токен"
// @Failure 400 {object} errs.AppError "Некорректный запрос"
// @Failure 401 {object} errs.AppError "Невалидный Google токен"
// @Failure 500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	user, err := h.googleAuthService.GoogleAuth(req.Token)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	if user.ProfileComplete {
		session := sessions.Default(c)
		session.Set(constants.ContextUserID, user.UserID)
		if err = session.Save(); err != nil {
			errs.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, dto.LoginResponse{
			Message:   constants.LoggedIn,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
		return
	}

	tempToken, err := h.tokenService.GenerateIDToken(user.UserID)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.RegistrationIncompleteResponse{
		Message: constants.UncompletedRegistration,
		Token:   tempToken,
	})
}

// CompleteSignUp завершает регистрацию пользователя, начавшего вход через Google.
//
// @Summary Завершение регистрации
// @Description Заполняет недостающие поля профиля после входа через Google.
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "ID токен в формате Bearer <token>"
// @Param first_name formData string true "Имя"
// @Param last_name formData string true "Фамилия"
// @Param second_name formData string false "Отчество"
// @Param birthday formData string true "Дата рождения (в формате 02.01.2006)"
// @Param password formData string true "Пароль"
// @Param region_id formData int true "ID региона"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} errs.AppError "Некорректный формат даты или другие ошибки валидации"
// @Failure 401 {object} errs.AppError "Отсутствует или невалидный/истёкший токен авторизации"
// @Failure 403 {object} errs.AppError "Регистрация уже завершена"
// @Failure 404 {object} errs.AppError "Регион не найден"
// @Failure 500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Security ApiToken
// @Router /auth/complete-sign-up [post]
func (h *AuthHandler) CompleteSignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	userID := c.MustGet(constants.ContextUserID).(uint)

	user, err := h.googleAuthService.CompleteProfile(userID, &req)
	if err != nil {
		errs.HandleError(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set(constants.ContextUserID, userID)
	if err = session.Save(); err != nil {
		errs.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message:   constants.LoggedIn,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	})
}
