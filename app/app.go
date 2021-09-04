package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/muchlist/berita_acara/configs"
	"github.com/muchlist/berita_acara/db"
	"github.com/muchlist/berita_acara/utils/logger"
	"log"
	"os"
	"os/signal"
)

func RunApp(){
	// Init config, logger dan db
	configs.Init()
	logger.Init()
	db.Init()
	defer db.Close()

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

	// ...

	// blocking and listen for fiber
	if err := app.Listen(":3500"); err != nil {
		logger.Error("error fiber listen", err)
		log.Panic()
	}

	// cleanup app
	fmt.Println("Running cleanup tasks...")
}