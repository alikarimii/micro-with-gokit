package rest

import (
	"github.com/go-kit/kit/log"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName       string `envconfig:"SERVICE_NAME"`
	TracerExporterUrl string `envconfig:"TRACER_EXPORTER_URL"`
	TracerID          int    `envconfig:"TRACER_ID"`
	Environment       string `envconfig:"ENVIRONMENT"`
	DatebaseName      string `envconfig:"DATABSE_NAME"`
	DatebaseUrl       string `envconfig:"DATABASE_URL"`
	HttpPort          string `envconfig:"HTTP_PORT"`
	HttpDialTimeout   int    `envconfig:"HTTP_DIAL_TIMEOUT"`
	MigrationsPath    string `envconfig:"DATABASE_MIGRATION_PATH"`
	LogLevel          string `envconfig:"LOG_LEVEL"`
}

func MustBuildConfigFromEnv(logger log.Logger) *Config {
	conf := Config{}
	mapTo("CONF", &conf, logger)

	if conf.ServiceName == "" {
		logger.Log("env", "env load faild: ServiceName", conf.ServiceName)
		panic("")
	}

	if conf.TracerExporterUrl == "" {
		logger.Log("env", "env load faild: TracerExporterUrl", conf.TracerExporterUrl)
		panic("")
	}
	if conf.DatebaseName == "" {
		logger.Log("env", "env load faild: DatebaseName", conf.DatebaseName)
		panic("")
	}
	if conf.DatebaseUrl == "" {
		logger.Log("env", "env load faild: DatebaseUrl", conf.DatebaseUrl)
		panic("")
	}
	if conf.HttpPort == "" {
		logger.Log("env", "env load faild: HttpPort", conf.HttpPort)
		panic("")
	}
	if conf.MigrationsPath == "" {
		logger.Log("env", "env load faild: MigrationPath", conf.MigrationsPath)
		panic("")
	}
	if conf.LogLevel == "" {
		logger.Log("env", "env load faild: LogLevel", conf.LogLevel)
		panic("")
	}
	if conf.HttpDialTimeout == 0 {
		logger.Log("env", "env load faild: HttpDialTimeout", conf.HttpDialTimeout)
		panic("")
	}
	if conf.Environment == "" {
		logger.Log("env", "env load faild: Environment", conf.Environment)
		panic("")
	}
	if conf.TracerID == 0 {
		logger.Log("env", "env load faild: TracerID", conf.TracerID)
		panic("")
	}

	return &conf
}

// mapTo map section
func mapTo(env string, schema interface{}, logger log.Logger) {
	e := envconfig.Process(env, schema)
	if e != nil {
		logger.Log("env", "can't load env: ", e)
		panic("")
	}
}
