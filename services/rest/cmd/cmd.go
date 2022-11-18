package cmd

import (
	"context"
	"os"
	"time"

	"github.com/alikarimii/micro-with-gokit/services/rest"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel"
)

func RunHttp() *rest.Service {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	config := rest.MustBuildConfigFromEnv(logger)
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(config.HttpDialTimeout)*time.Second)
	// Set up logger with level filter.
	var opt level.Option
	switch config.LogLevel {
	case "info":
		opt = level.AllowInfo()
	case "debug":
		opt = level.AllowDebug()
	default:
		opt = level.AllowNone()
	}
	logger = level.NewFilter(logger, opt)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	//
	// db
	dbConn := rest.MustInitPostgresDB(config, logger)
	//tracer
	tp := rest.MustBuildTracer(config)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	tr := tp.Tracer("rest.Service")

	ctx, span := tr.Start(ctx, "RunHttp")
	defer span.End()

	// di
	di := rest.MustBuildDIContainer(
		config,
		logger,
		rest.WithPostgresConnection(dbConn),
		rest.WithMetrics(),
		rest.WithJaeger(tp),
	)

	exitFn := func() {
		cancelFn()
		tp.Shutdown(ctx)
		os.Exit(1)
	}
	return rest.InitService(ctx, config, logger, exitFn, di)
}
