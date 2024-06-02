package config

import (
	"encoding/json"
	"os"
)

type Appsettings struct {
	Config   Config
	Logging  Logging
	Database Database
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type Logging struct {
	EnableConsole bool
	LogFilePath   string
}

type Database struct {
	User string
	Pass string
	Host string
	Db   string
	Port uint16
}

var appsettings Appsettings

func Init() {
	bytes, err := os.ReadFile("appsettings.json")
	if err != nil {
		panic(err.Error() + "\nCOULD NOT OPEN appsettings.json")
	}

	var result *Appsettings
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		panic(err.Error() + "\nCOULD NOT READ appsettings.json")
	}

	appsettings = *result
}

func GetAppsettings() Appsettings {
	return appsettings
}
