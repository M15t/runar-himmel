package main

import (
	"embed"
	"fmt"
	"runar-himmel/config"
	"runar-himmel/internal/api/auth"
	"runar-himmel/internal/api/root"
	"runar-himmel/internal/db"
	"runar-himmel/internal/rbac"
	"runar-himmel/internal/repo"

	"runar-himmel/pkg/server"
	"runar-himmel/pkg/server/middleware/jwt"
	"runar-himmel/pkg/server/middleware/secure"
	"runar-himmel/pkg/util/crypter"

	"github.com/labstack/echo/v4"
)

// To embed SwaggerUI into api server using go:build tag
var (
	enableSwagger = false
	swaggerui     embed.FS
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

	fmt.Println(db)

	// Initialize HTTP server
	e := server.New(&server.Config{
		Port:              cfg.Server.Port,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		AllowOrigins:      cfg.Server.AllowOrigins,
		Debug:             cfg.General.Debug,
	})

	// custom api context
	// e.Use(api.ContextMiddleware())

	if enableSwagger {
		// Static page for SwaggerUI
		e.GET("/swagger-ui*", echo.StaticDirectoryHandler(echo.MustSubFS(swaggerui, "swagger-ui"), false), secure.DisableCache())
	}

	// Initialize core services
	crypterSvc := crypter.New()
	repoSvc := repo.New(db)
	rbacSvc := rbac.New(cfg.General.Debug)
	jwtSvc := jwt.New(cfg.JWT.Algorithm, cfg.JWT.Secret, cfg.JWT.DurationAccessToken, cfg.JWT.DurationRefreshToken)

	fmt.Println(crypterSvc, rbacSvc, jwtSvc, repoSvc)

	// Initialize services
	authSvc := auth.New(repoSvc, jwtSvc, crypterSvc)

	// Initialize root API
	root.NewHTTP(e)

	auth.NewHTTP(authSvc, e.Group("/auth"))

	// ctx := context.Context(context.Background())
	// newUser := &types.User{
	// 	FirstName: "Runar",
	// 	LastName:  "Himmel",
	// 	Email:     "rn@runar.sky",
	// }

	// rec := &types.User{}
	// if err := repoSvc.User.GDB.Take(rec, `email = ?`, `loki@runar-himmel.sky`).Error; err != nil {
	// 	fmt.Println("====== err", err)
	// }

	// if err := repoSvc.User.Read(ctx, rec, `email = ?`, `loki@runar-himmel.sky`); err != nil {
	// 	fmt.Println("====== err", err)
	// }

	// fmt.Println("====== result", rec)

	server.Start(e, config.IsLambda())
}
