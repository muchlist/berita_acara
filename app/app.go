package app

import (
	"github.com/muchlist/berita_acara/configs"
	"github.com/muchlist/berita_acara/db"
	"github.com/muchlist/berita_acara/utils/logger"
)

func RunApp(){
	// Init config, logger and db
	configs.Init()
	logger.Init()
	db.Init()
	defer db.Close()



}