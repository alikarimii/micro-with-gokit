package rest

import (
	"database/sql"

	"github.com/alikarimii/micro-with-gokit/internal/application"
	"github.com/alikarimii/micro-with-gokit/internal/endpoints"
	myhttp "github.com/alikarimii/micro-with-gokit/internal/infrastructure/http"
	psql "github.com/alikarimii/micro-with-gokit/internal/infrastructure/postgres"
	"github.com/alikarimii/micro-with-gokit/internal/infrastructure/postgres/database"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type DIOption func(container *DIContainer) error

type DIContainer struct {
	config *Config
	logger log.Logger
	infra  struct {
		postgres    *sql.DB
		httpHandler *mux.Router
		endpoint    *endpoints.Endpoints
	}
	app struct {
		studentService application.StudentService
		repo           application.StudentRepository
	}
	trace struct {
		duration       metrics.Histogram
		counter        metrics.Counter
		jaegerExporter *tracesdk.TracerProvider
	}
}

func (container *DIContainer) init() {
	_ = container.GetStudentService()
	_ = container.GetPostgresStudentRepository()
	_ = container.GetEndpoints()
	_ = container.GetHttpHandler()
}

func (container *DIContainer) GetStudentService() application.StudentService {
	if container.app.studentService == nil {
		container.app.studentService = application.NewWithLogMiddlware(
			container.logger,
			container.GetPostgresStudentRepository(),
		)
	}
	return container.app.studentService
}
func (container *DIContainer) GetPostgresStudentRepository() application.StudentRepository {
	if container.app.repo == nil {
		container.app.repo = psql.NewRepo(
			container.infra.postgres,
			container.logger,
		)
	}
	return container.app.repo
}
func (container *DIContainer) GetEndpoints() *endpoints.Endpoints {
	if container.infra.endpoint == nil {
		container.infra.endpoint = endpoints.New(
			container.GetStudentService(),
			container.logger,
			container.trace.duration,
		)
	}
	return container.infra.endpoint
}
func (container *DIContainer) GetHttpHandler() *mux.Router {
	if container.infra.httpHandler == nil {
		container.infra.httpHandler = myhttp.NewHTTPHandler(
			container.GetEndpoints(),
			container.logger,
		)
	}
	return container.infra.httpHandler
}

func WithMetrics() DIOption {
	return func(container *DIContainer) error {
		// Endpoint-level metrics.
		container.trace.duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "studentService",
			Subsystem: "rest",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
		return nil
	}
}
func WithJaeger(tp *tracesdk.TracerProvider) DIOption {
	return func(container *DIContainer) error {
		container.trace.jaegerExporter = tp
		return nil
	}
}
func WithPostgresConnection(dbConn *sql.DB) DIOption {
	return func(container *DIContainer) error {

		container.infra.postgres = dbConn
		return nil
	}
}

func MustBuildTracer(config *Config) *tracesdk.TracerProvider {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.TracerExporterUrl)))
	if err != nil {
		panic(err)
	}
	return tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			attribute.String("environment", config.Environment),
			attribute.Int("ID", config.TracerID),
		)),
	)
}
func MustBuildDIContainer(config *Config, logger log.Logger, opts ...DIOption) *DIContainer {
	container := &DIContainer{
		config: config,
		logger: logger,
	}

	for _, v := range opts {
		if err := v(container); err != nil {
			logger.Log("build", "mustBuildDIContainer: ", err)
			panic(err)
		}
	}
	container.init()

	return container
}

func MustInitPostgresDB(config *Config, logger log.Logger) *sql.DB {
	var err error

	logger.Log("build db", "bootstrapPostgresDB: opening Postgres DB connection ...")

	postgresDBConn, err := sql.Open("postgres", config.DatebaseUrl)
	if err != nil {
		logger.Log("build db", "bootstrapPostgresDB: failed to open Postgres DB connection: ", err)
		panic("")
	}

	err = postgresDBConn.Ping()
	if err != nil {
		logger.Log("build db", "bootstrapPostgresDB: failed to connect to Postgres DB: ", err)
		panic(err)
	}

	/***/
	logger.Log("build db", "bootstrapPostgresDB: running DB migrations for student ...")

	migrateStudent, err := database.NewMigrator(postgresDBConn, config.MigrationsPath)
	if err != nil {
		logger.Log("build db", "bootstrapPostgresDB: failed to create DB migrator for student:")
		panic(err)
	}
	// migrateStudent.WithLogger(logger)
	err = migrateStudent.Up()
	if err != nil {
		logger.Log("build db", "bootstrapPostgresDB: failed to run DB migrations for student: ")
		panic(err)
	}

	return postgresDBConn
}
