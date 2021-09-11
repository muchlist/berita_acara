package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/berita_acara/configs"
	"github.com/muchlist/berita_acara/utils/logger"
)

var (
	DB *pgxpool.Pool
)

// Init menginisiasi database pool
func Init() {
	cfg := configs.Config

	// databaseUrl := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUSER, cfg.DBPASS, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)

	var err error
	DB, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		logger.Error("tidak dapat terhubung ke database", err)
		panic("Unable to connect to database")
	}

	logger.Info("terkoneksi ke database")
}

func Close() {
	DB.Close()
	logger.Info("koneksi ke database ditutup")
}
