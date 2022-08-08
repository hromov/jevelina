package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment int

const (
	EnvLocal = iota
	EnvStaging
	EnvProd
)

type Config struct {
	Env        Environment
	Dsn        string
	BucketName string
}

var defaultConfig = Config{
	Env:        EnvLocal,
	BucketName: "jevelina",
}

func Get() Config {
	const envConfPrefix = "DUNGEON"

	loadedConfig := defaultConfig
	envconfig.MustProcess(envConfPrefix, &loadedConfig)
	if loadedConfig.Env != EnvLocal {
		return loadedConfig
	} else if err := godotenv.Load(); err == nil {
		envconfig.MustProcess(envConfPrefix, &loadedConfig)
	}

	return loadedConfig
}
