package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/berita_acara/configs"
)

var (
	DB *pgxpool.Pool
)

// Init menginisiasi database pool
// responsenya digunakan untuk memutus koneksi apabila main program dihentikan
func Init() *pgxpool.Pool {
	cfg := configs.Config

	// databaseUrl := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUSER, cfg.DBPASS, cfg.DBHOST, cfg.DBHOST, cfg.DBNAME)

	var err error
	DB, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	fmt.Println("Connected!")

	return DB
}