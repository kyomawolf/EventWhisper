package configuration

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port         int
	LogLevel     string
	BasePath     string
	DBConnection string
	DatabaseName string
}

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	return v
}
func getenvInt(key string, defaultValue int) int {
	s := getenvStr(key, "")
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return v
}

func LoadConfig() (*Config, error) {
	return &Config{
		Port:         getenvInt("PORT", 8080),
		LogLevel:     getenvStr("LOG_LEVEL", "info"),
		BasePath:     "",
		DBConnection: getenvStr("MONGO_CONNECTION", "mongodb://root:example@localhost:27017"),
		DatabaseName: "eventwhisper",
	}, nil
}
