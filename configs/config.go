package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configuration struct {
	DBUSER    string
	DBPASS    string
	DBHOST    string
	DBPORT    string
	DBNAME    string
	LOGLEVEL  string
	SECRETKEY string
}

var (
	Config *Configuration
)

func SetupConfiguration() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("File .env tidak ditemukan")
	}

	Config = new(Configuration)

	Config.DBUSER = os.Getenv("BA_DB_USER")
	Config.DBPASS = os.Getenv("BA_DB_PASS")
	Config.DBHOST = os.Getenv("BA_DB_HOST")
	Config.DBPORT = os.Getenv("BA_DB_PORT")
	Config.DBNAME = os.Getenv("BA_DB_NAME")
	Config.LOGLEVEL = os.Getenv("BA_LOG_LEVEL")
	Config.SECRETKEY = os.Getenv("BA_SECRET_KEY")
}
