package utils

import "os"

type Config struct {
	DbUser string
	DbPassword string
	DbName string
	TgToken string
	Host string
	Port string
}

func ReadConfig() Config {
	var config Config
	config.DbUser = os.Getenv("DB_USER")
	config.DbPassword = os.Getenv("DB_PASSWORD")
	config.DbName = os.Getenv("DB_NAME")
	config.TgToken = os.Getenv("TG_TOKEN")
	config.Host = os.Getenv("DB_HOST")
	config.Port = os.Getenv("DB_PORT")

	return config
}