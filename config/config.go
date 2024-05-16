package config

import (
	"encoding/json"
	"os"
)

type Appsettings struct {
	Config Config
}

type Config struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

func New() Appsettings {
	bytes, err := os.ReadFile("appsettings.json")
	if err != nil {
		panic(err.Error() + "\nCOULD NOT OPEN appsettings.json")
	}

	var result *Appsettings
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		panic(err.Error() + "\nCOULD NOT READ appsettings.json")
	}

	return *result
}
