package config

import "os"

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	Port            string
	TraqAccessToken string
	BotName         string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:          getEnv("NS_MARIADB_HOSTNAME", "localhost"),
		DBPort:          getEnv("NS_MARIADB_PORT", "3306"),
		DBUser:          getEnv("NS_MARIADB_USER", "user"),
		DBPassword:      getEnv("NS_MARIADB_PASSWORD", "password"),
		DBName:          getEnv("NS_MARIADB_DATABASE", "template_db"),
		Port:            getEnv("PORT", "8080"),
		TraqAccessToken: getEnv("TRAQ_ACCESS_TOKEN", ""),
		BotName:         getEnv("BOT_NAME", "reminder-bot"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
