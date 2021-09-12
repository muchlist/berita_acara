package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/configs"
	"github.com/muchlist/berita_acara/db"
	_ "github.com/muchlist/berita_acara/docs"
	"github.com/muchlist/berita_acara/utils/logger"
	"github.com/muchlist/berita_acara/utils/mjwt"
	"log"
	"os"
	"os/signal"
)

// RunApp
// @title Berita Acara API
// @version 1.0
// @description Berita Acara Api
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email whois.muchlis@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
// @host localhost:3500
// @BasePath /api/v1
func RunApp() {
	// Init config, logger dan db
	configs.Init()
	logger.Init()
	db.Init()
	defer db.Close()
	mjwt.Init()

	// membuat fiber app
	app := fiber.New()

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	prepareEndPoint(app)

	// blocking and listen for fiber
	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		log.Panic()
	}

	// cleanup app
	fmt.Println("Running cleanup tasks...")
}
