package app

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/configs/roles"
	"github.com/muchlist/berita_acara/dao/userdao"
	"github.com/muchlist/berita_acara/db"
	"github.com/muchlist/berita_acara/handler"
	"github.com/muchlist/berita_acara/middle"
	"github.com/muchlist/berita_acara/services/userserv"
	"github.com/muchlist/berita_acara/utils/mcrypt"
	"github.com/muchlist/berita_acara/utils/mjwt"
)

func prepareEndPoint(app *fiber.App) {

	// Utils
	cryptoUtils := mcrypt.NewCrypto()
	jwt := mjwt.NewJwt()

	// User Domain
	userDao := userdao.New(db.DB)
	userService := userserv.NewUserService(userDao, cryptoUtils, jwt)
	userHandler := handler.NewUserHandler(userService)

	app.Get("/swagger/*", swagger.Handler) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://localhost.com/docs/swagger.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	// url mapping
	api := app.Group("/api/v1")

	//USER
	api.Get("/users/:id", userHandler.Get)
	api.Get("/users", middle.NormalAuth(), userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middle.NormalAuth(), userHandler.GetProfile)
	api.Post("/register", middle.NormalAuth(roles.RoleAdmin), userHandler.Register)
	api.Put("/users/:id", middle.NormalAuth(roles.RoleAdmin), userHandler.Edit)
	api.Delete("/users/:id", middle.NormalAuth(roles.RoleAdmin), userHandler.Delete)
}
