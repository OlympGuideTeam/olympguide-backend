package handler

import (
	"api/dto"
	"api/service"
	"api/utils/constants"
	"api/utils/errs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService         service.IAuthService
	externalAuthService service.IExternalAuthService
	tokenService        service.ITokenService
}

func NewAuthHandler(
	authService service.IAuthService,
	externalAuthService service.IExternalAuthService,
	tokenService service.ITokenService,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		externalAuthService: externalAuthService,
		tokenService:        tokenService,
	}
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

// SignUp godoc
// @Summary      Регистрация пользователя по email
// @Description  Создаёт нового пользователя по переданному email и паролю. Должен быть подтверждён Email и получен токен на этапе verify_code. После успешной регистрации создаётся сессия.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.EmailSignUpRequest true "Данные для регистрации"
// @Success      201 {object} dto.LoginResponse "Пользователь успешно зарегистрирован, сессия создана"
// @Failure      400 {object} errs.AppError "Неверный формат запроса"
// @Failure      409 {object} errs.AppError "Пользователь с таким email уже существует"
// @Failure      500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Security     ApiToken
// @Router       /auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var request dto.EmailSignUpRequest
	if err := c.ShouldBind(&request); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}
	request.Email = c.MustGet(constants.ContextEmail).(string)

	user, err := h.authService.SignUp(&request)
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

	c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: constants.SignedUp,
	})
}

// Login godoc
// @Summary      Вход пользователя
// @Description  Авторизует пользователя по email и паролю. В случае успеха создаётся сессия, возвращаются имя и фамилия пользователя.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Данные для входа"
// @Success      200 {object} dto.LoginResponse "Успешный вход, сессия создана"
// @Failure      400 {object} errs.AppError "Неверный формат запроса"
// @Failure      404 {object} errs.AppError "Пользователя с таким email не существует"
// @Failure      500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Security     SessionCookie
// @Router       /auth/login [post]
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

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: constants.LoggedIn,
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
// @Summary Вход через Google
// @Description При успешном входе устанавливается сессия.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ExternalAuthRequest true "Токен Google ID"
// @Success 200 {object} dto.LoginResponse "Регистрация завершена — пользователь вошёл"
// @Failure 400 {object} errs.AppError "Некорректный запрос"
// @Failure 401 {object} errs.AppError "Невалидный Google токен"
// @Failure 500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.ExternalAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	user, err := h.externalAuthService.GoogleAuth(req.Token)
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

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: constants.LoggedIn,
	})
	return
}

// AppleLogin выполняет вход через Apple.
// @Summary Вход через Apple
// @Description При успешном входе устанавливается сессия.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ExternalAuthRequest true "Токен Apple ID"
// @Success 200 {object} dto.LoginResponse "Регистрация завершена — пользователь вошёл"
// @Failure 400 {object} errs.AppError "Некорректный запрос"
// @Failure 401 {object} errs.AppError "Невалидный Apple токен"
// @Failure 500 {object} errs.AppError "Внутренняя ошибка сервера"
// @Router /auth/apple [post]
func (h *AuthHandler) AppleLogin(c *gin.Context) {
	var req dto.ExternalAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.HandleError(c, errs.InvalidRequest)
		return
	}

	user, err := h.externalAuthService.AppleAuth(req.Token)
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

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: constants.LoggedIn,
	})
	return
}
