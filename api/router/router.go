package router

import (
	"api/handler"
	"api/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler   *handler.AuthHandler
	univerHandler *handler.UniverHandler
	fieldHandler  *handler.FieldHandler
	olympHandler  *handler.OlympHandler
	metaHandler   *handler.MetaHandler
}

func NewRouter(auth *handler.AuthHandler, univer *handler.UniverHandler,
	field *handler.FieldHandler, olymp *handler.OlympHandler, meta *handler.MetaHandler) *Router {
	return &Router{
		authHandler:   auth,
		univerHandler: univer,
		fieldHandler:  field,
		olympHandler:  olymp,
		metaHandler:   meta,
	}
}

func (rt *Router) SetupRoutes(r *gin.Engine) {
	r.Use(middleware.SessionMiddleware())
	r.Use(middleware.ValidateID())

	rt.setupAuthRoutes(r)
	rt.setupUniverRoutes(r)
	rt.setupUserRoutes(r)
	rt.setupFieldRoutes(r)
	rt.setupOlympRoutes(r)
	rt.setupMetaRoutes(r)
}

func (rt *Router) setupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/send-code", rt.authHandler.SendCode)
	authGroup.POST("/verify-code", rt.authHandler.VerifyCode)
	authGroup.POST("/sign-up", rt.authHandler.SignUp)
	authGroup.POST("/login", rt.authHandler.Login)
	authGroup.POST("/logout", rt.authHandler.Logout)
}

func (rt *Router) setupUniverRoutes(r *gin.Engine) {
	r.GET("/universities", rt.univerHandler.GetUnivers)
	university := r.Group("/university")
	{
		secured := university.Group("/", middleware.AuthMiddleware())
		{
			secured.POST("/", rt.univerHandler.NewUniver)
			secured.PUT("/:id", rt.univerHandler.UpdateUniver)
			secured.DELETE("/:id", rt.univerHandler.DeleteUniver)
		}
	}
}

func (rt *Router) setupFieldRoutes(r *gin.Engine) {
	r.GET("/fields", rt.fieldHandler.GetGroups)
	r.GET("/field/:id", rt.fieldHandler.GetField)
}

func (rt *Router) setupOlympRoutes(r *gin.Engine) {
	r.GET("/olympiads", rt.olympHandler.GetOlympiads)
}

func (rt *Router) setupUserRoutes(r *gin.Engine) {
	user := r.Group("/user", middleware.AuthMiddleware())
	{
		//userGroup.GET("/region", user.GetRegion)
		favourite := user.Group("/favourite")
		{
			favourite.GET("/universities", rt.univerHandler.GetLikedUnivers)
			favourite.POST("/university/:id", rt.univerHandler.LikeUniver)
			favourite.DELETE("/university/:id", rt.univerHandler.DislikeUniver)

			favourite.GET("/olympiad", rt.olympHandler.GetLikedOlympiads)
			favourite.POST("/olympiad/:id", rt.olympHandler.LikeOlymp)
			favourite.DELETE("/olympiad/:id", rt.olympHandler.DislikeOlymp)
		}
	}
}

func (rt *Router) setupMetaRoutes(r *gin.Engine) {
	meta := r.Group("/meta")
	meta.GET("/regions", rt.metaHandler.GetRegions)
	meta.GET("/university-regions", rt.metaHandler.GetUniversityRegions)
	meta.GET("/olympiad-profiles", rt.metaHandler.GetOlympiadProfiles)
}
