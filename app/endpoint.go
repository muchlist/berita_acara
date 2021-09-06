package app

import (
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

	// url mapping
	api := app.Group("/api/v1")

	//USER
	api.Get("/users/:id", userHandler.Get)
	api.Get("/users", userHandler.Find)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.RefreshToken)
	api.Get("/profile", middle.NormalAuth(), userHandler.GetProfile)
	api.Post("/register-force", userHandler.Register)                               // <- seharusnya gunakan middleware agar hanya admin yang bisa meregistrasi
	api.Post("/register", middle.NormalAuth(roles.RoleAdmin), userHandler.Register) // <- hanya admin yang bisa meregistrasi
	api.Put("/users/:id", middle.NormalAuth(roles.RoleAdmin), userHandler.Edit)
	api.Delete("/users/:id", middle.NormalAuth(roles.RoleAdmin), userHandler.Delete)
}
