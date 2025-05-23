package router

import (
	_ "api/docs"
	"api/handler"
	"api/middleware"
	"api/utils/role"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

type Router struct {
	engine   *gin.Engine
	api      *gin.RouterGroup
	handlers *handler.Handlers
	mw       *middleware.Mw
	store    sessions.Store
}

func NewRouter(handlers *handler.Handlers, mw *middleware.Mw, store sessions.Store) *Router {
	engine := gin.Default()
	apiGroup := engine.Group("/api/v1")

	router := &Router{
		engine:   engine,
		api:      apiGroup,
		handlers: handlers,
		mw:       mw,
		store:    store,
	}

	router.setupRoutes()
	return router
}

func (rt *Router) Run(port int) {
	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("Server listening on %s", serverAddress)
	if err := rt.engine.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (rt *Router) setupRoutes() {
	rt.api.Use(sessions.Sessions("session", rt.store))
	rt.api.Use(rt.mw.PrometheusMetrics())
	rt.api.Use(rt.mw.SessionMiddleware())
	rt.api.Use(rt.mw.ValidateNumericParams())

	rt.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	rt.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rt.setupAuthRoutes()
	rt.setupUniverRoutes()
	rt.setupUserRoutes()
	rt.setupFieldRoutes()
	rt.setupOlympRoutes()
	rt.setupMetaRoutes()
	rt.setupFacultyRoutes()
	rt.setupProgramRoutes()
	rt.setupBenefitRoutes()
	rt.setupServiceDiplomaRoutes()
}

func (rt *Router) setupAuthRoutes() {
	authGroup := rt.api.Group("/auth")
	authGroup.POST("/send-code", rt.handlers.Auth.SendCode)
	authGroup.POST("/verify-code", rt.handlers.Auth.VerifyCode)
	authGroup.POST("/sign-up", rt.mw.AlreadyLoginMiddleware(), rt.mw.EmailTokenMiddleware(), rt.handlers.Auth.SignUp)
	authGroup.POST("/login", rt.mw.AlreadyLoginMiddleware(), rt.handlers.Auth.Login)
	authGroup.POST("/google", rt.mw.AlreadyLoginMiddleware(), rt.handlers.Auth.GoogleLogin)
	authGroup.POST("/apple", rt.mw.AlreadyLoginMiddleware(), rt.handlers.Auth.AppleLogin)
	authGroup.POST("/logout", rt.mw.UserMiddleware(), rt.handlers.Auth.Logout)
	authGroup.GET("/check-session", rt.handlers.Auth.CheckSession)
}

func (rt *Router) setupUniverRoutes() {
	rt.api.GET("/universities", rt.handlers.Univer.GetUnivers)
	university := rt.api.Group("/university")
	{
		university.POST("/", rt.mw.RolesMiddleware(role.Founder, role.Admin, role.DataLoaderService), rt.handlers.Univer.NewUniver)
		universityWithID := university.Group("/:id")
		{
			universityWithID.GET("/", rt.handlers.Univer.GetUniver)
			universityWithID.POST("/logo",
				rt.mw.RolesMiddleware(role.Founder, role.Admin),
				rt.handlers.Univer.UploadUniverLogo)
			universityWithID.DELETE("/logo",
				rt.mw.RolesMiddleware(role.Founder, role.Admin),
				rt.handlers.Univer.DeleteLogo)
			universityWithID.GET("/faculties", rt.handlers.Faculty.GetFaculties)
			universityWithID.GET("/programs/by-faculty", rt.handlers.Program.GetUniverProgramsWithFaculty)
			universityWithID.GET("/programs/by-field", rt.handlers.Program.GetUniverProgramsWithGroup)
		}

		universityWithID.Use(rt.mw.UniversityIdSetter(),
			rt.mw.RolesMiddleware(role.Founder, role.Admin, role.DataLoaderService, role.UniverEditor))
		{
			universityWithID.PUT("/", rt.handlers.Univer.UpdateUniver)
			universityWithID.DELETE("/", rt.handlers.Univer.DeleteUniver)
		}
	}
}

func (rt *Router) setupFieldRoutes() {
	rt.api.GET("/fields", rt.handlers.Field.GetGroups)
	rt.api.GET("/field/:id", rt.handlers.Field.GetField)
	rt.api.GET("/field/:id/programs", rt.handlers.Program.GetProgramsByField)
}

func (rt *Router) setupOlympRoutes() {
	rt.api.GET("/olympiads", rt.handlers.Olymp.GetOlympiads)
	olympiad := rt.api.Group("/olympiad")
	olympiadWithID := olympiad.Group("/:id")
	{
		olympiadWithID.GET("/", rt.handlers.Olymp.GetOlympiad)
		olympiadWithID.GET("/benefits", rt.handlers.Benefit.GetBenefitsByOlympiad)
		olympiadWithID.GET("/universities", rt.handlers.Univer.GetBenefitByOlympUnivers)
	}
}

func (rt *Router) setupUserRoutes() {
	user := rt.api.Group("/user", rt.mw.UserMiddleware())
	{
		user.DELETE("", rt.handlers.User.DeleteUser)
		user.GET("/data", rt.handlers.User.GetUserData)
		user.POST("/update", rt.handlers.User.UpdateUser)
		user.POST("/update-password", rt.mw.EmailTokenMiddleware(), rt.handlers.User.UpdatePassword)
		diplomas := user.Group("/diplomas")
		{
			diplomas.GET("/", rt.handlers.Diploma.GetUserDiplomas)
			diplomas.GET("/universities", rt.handlers.Univer.GetUserDiplomasUnivers)
			diplomas.GET("/benefits", rt.handlers.Benefit.GetUserBenefits)
			diplomas.POST("/sync", rt.handlers.Diploma.SyncUserDiplomas)
		}
		diploma := user.Group("/diploma")
		{
			diploma.POST("/", rt.handlers.Diploma.NewDiplomaByUser)
			diploma.GET("/:id/universities", rt.handlers.Univer.GetDiplomaUnivers)
			diploma.GET("/:id/benefits", rt.handlers.Benefit.GetBenefitsByDiploma)
			diploma.DELETE("/:id", rt.handlers.Diploma.DeleteDiploma)
		}

		favourite := user.Group("/favourite")
		{
			favourite.GET("/programs", rt.handlers.Program.GetLikedPrograms)
			favourite.POST("/program/:id", rt.handlers.Program.LikeProgram)
			favourite.DELETE("/program/:id", rt.handlers.Program.DislikeProgram)

			favourite.GET("/universities", rt.handlers.Univer.GetLikedUnivers)
			favourite.POST("/university/:id", rt.handlers.Univer.LikeUniver)
			favourite.DELETE("/university/:id", rt.handlers.Univer.DislikeUniver)

			favourite.GET("/olympiads", rt.handlers.Olymp.GetLikedOlympiads)
			favourite.POST("/olympiad/:id", rt.handlers.Olymp.LikeOlymp)
			favourite.DELETE("/olympiad/:id", rt.handlers.Olymp.DislikeOlymp)
		}
	}
}

func (rt *Router) setupMetaRoutes() {
	meta := rt.api.Group("/meta")
	meta.GET("/regions", rt.handlers.Meta.GetRegions)
	meta.GET("/university-regions", rt.handlers.Meta.GetUniversityRegions)
	meta.GET("/olympiad-profiles", rt.handlers.Meta.GetOlympiadProfiles)
	meta.GET("/subjects", rt.handlers.Meta.GetSubjects)
}

func (rt *Router) setupFacultyRoutes() {
	faculty := rt.api.Group("/faculty")
	faculty.Use(rt.mw.RolesMiddleware(role.Founder, role.Admin, role.DataLoaderService))
	{
		faculty.POST("/", rt.handlers.Faculty.NewFaculty)

		facultyWithID := faculty.Group("/:id")
		{
			facultyWithID.PUT("/", rt.handlers.Faculty.UpdateFaculty)
			facultyWithID.DELETE("/", rt.handlers.Faculty.DeleteFaculty)
			facultyWithID.GET("/programs", rt.mw.NoMiddleware(), rt.handlers.Program.GetProgramsByFaculty)
		}
	}
}

func (rt *Router) setupProgramRoutes() {
	program := rt.api.Group("/program")
	program.POST("/", rt.handlers.Program.NewProgram)
	programWithID := program.Group("/:id")
	{
		programWithID.GET("/", rt.handlers.Program.GetProgram)
		programWithID.GET("/benefits", rt.handlers.Benefit.GetBenefitsByProgram)
	}
}

func (rt *Router) setupBenefitRoutes() {
	benefit := rt.api.Group("/benefit")
	benefit.Use(rt.mw.RolesMiddleware(role.DataLoaderService, role.Admin, role.Founder))
	benefit.POST("/", rt.handlers.Benefit.NewBenefit)
	benefit.DELETE("/:id", rt.handlers.Benefit.DeleteBenefit)
}

func (rt *Router) setupServiceDiplomaRoutes() {
	rt.api.POST("/service/diploma", rt.mw.RolesMiddleware(role.DataLoaderService), rt.handlers.Diploma.NewDiplomaByService)
}
