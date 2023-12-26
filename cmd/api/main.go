package main

import (
	"runar-himmel/config"
	"runar-himmel/internal/api/admin/session"
	"runar-himmel/internal/api/admin/user"
	"runar-himmel/internal/api/auth"
	"runar-himmel/internal/api/root"
	"runar-himmel/internal/db"
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/repo"

	"runar-himmel/pkg/server"
	"runar-himmel/pkg/server/middleware/jwt"
	"runar-himmel/pkg/util/crypter"

	contextutil "runar-himmel/internal/api/context"

	_ "runar-himmel/cmd/api/docs" // Swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Runar Himmel server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8081
//	@BasePath	/v1

// To embed SwaggerUI into api server using go:build tag
var (
	enableSwagger = true
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := config.LoadAll()
	checkErr(err)

	db, sqldb, err := db.New(cfg.DB)
	checkErr(err)
	defer sqldb.Close()

	// Initialize HTTP server
	e := server.New(&server.Config{
		Port:              cfg.Server.Port,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		AllowOrigins:      cfg.Server.AllowOrigins,
		Debug:             cfg.General.Debug,
	})

	if enableSwagger {
		// Static page for SwaggerUI (go-swagger)
		// e.GET("/swagger-ui*", echo.StaticDirectoryHandler(echo.MustSubFS(swaggerui, "swagger-ui"), false), secure.DisableCache())

		// Static page for SwaggerUI (swaggo)
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Initialize core services
	crypterSvc := crypter.New()
	repoSvc := repo.New(db)
	rbacSvc := rbac.New(cfg.General.Debug)
	jwtSvc := jwt.New(cfg.JWT.Algorithm, cfg.JWT.Secret, cfg.JWT.DurationAccessToken, cfg.JWT.DurationRefreshToken)

	// Initialize services
	authSvc := auth.New(repoSvc, jwtSvc, crypterSvc)
	sessionSvc := session.New(repoSvc, rbacSvc)
	userSvc := user.New(repoSvc, rbacSvc, crypterSvc)

	// Initialize root API
	root.NewHTTP(e)

	auth.NewHTTP(authSvc, e.Group("/auth"))

	// Initialize admin APIs
	adminRouter := e.Group("/admin")
	adminRouter.Use(jwtSvc.MWFunc(), contextutil.MWContext())
	session.NewHTTP(sessionSvc, adminRouter.Group("/sessions"))
	user.NewHTTP(userSvc, adminRouter.Group("/users"))

	server.Start(e, config.IsLambda())
}
