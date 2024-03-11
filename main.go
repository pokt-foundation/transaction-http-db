package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgresdriver "github.com/pokt-foundation/transaction-db/postgres-driver"
	"github.com/pokt-foundation/transaction-http-db/batch"
	"github.com/pokt-foundation/transaction-http-db/metric"
	"github.com/pokt-foundation/transaction-http-db/router"
	"github.com/pokt-foundation/utils-go/environment"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Postgres DB vars - Required for all Envs.
	pgUser     = "PG_USER"
	pgPassword = "PG_PASSWORD"
	pgDatabase = "PG_DATABASE"
	// Local DB vars - Required for development/test Env.
	pgHost = "PG_HOST"
	pgPort = "PG_PORT"
	// CloudSQL DB vars - Required for production Env.
	dbInstanceConnectionName = "DB_INSTANCE_CONNECTION_NAME"
	privateIP                = "PRIVATE_IP"

	chanSize                      = "CHAN_SIZE"
	apiKeys                       = "API_KEYS"
	port                          = "PORT"
	maxRelayBatchSize             = "MAX_RELAY_BATCH_SIZE"
	maxRelayBatchDuration         = "MAX_RELAY_BATCH_DURATION"
	maxServiceRecordBatchSize     = "MAX_SERVICE_RECORD_BATCH_SIZE"
	maxServiceRecordBatchDuration = "MAX_SERVICE_RECORD_BATCH_DURATION"
	dbTimeout                     = "DB_TIMEOUT"
	debug                         = "DEBUG"

	defaultPort          = "8080"
	defaultBatchSize     = 1000
	defaultBatchDuration = 60
	defaultDBTimeout     = 60
	defaultDebug         = false
	defaultUseSSH        = false
	defaultChanSize      = 10000
)

type (
	options struct {
		// Required vars
		apiKeys                        map[string]bool
		pgUser, pgPassword, pgDatabase string
		// Local DB vars - Required for development/test Env.
		pgHost, pgPort string
		// CloudSQL DB vars - Required for production Env.
		dbInstanceConnectionName string
		privateIP                bool
		// Optional vars
		port                          string
		maxRelayBatchSize             int
		maxRelayBatchDuration         time.Duration
		maxServiceRecordBatchSize     int
		maxServiceRecordBatchDuration time.Duration
		dbTimeout                     time.Duration
		debug                         bool
		chanSize                      int
	}

	// DB config structs
	DBConfig interface {
		GetDriver(ctx context.Context) (driver *postgresdriver.PostgresDriver, cleanup func() error, err error)
	}
	cloudSQLConfig struct {
		options
	}
	testDBConfig struct {
		options
	}
)

func gatherOptions() options {
	return options{
		// Required vars
		apiKeys:    environment.MustGetStringMap(apiKeys, ","),
		pgUser:     environment.MustGetString(pgUser),
		pgPassword: environment.MustGetString(pgPassword),
		pgDatabase: environment.MustGetString(pgDatabase),
		// CloudSQL DB Config var
		dbInstanceConnectionName: environment.GetString(dbInstanceConnectionName, ""),
		privateIP:                environment.GetBool(privateIP, false),
		// Local DB Config vars
		pgHost: environment.GetString(pgHost, ""),
		pgPort: environment.GetString(pgPort, ""),
		// Optional vars
		port:                          environment.GetString(port, defaultPort),
		maxRelayBatchSize:             int(environment.GetInt64(maxRelayBatchSize, defaultBatchSize)),
		maxRelayBatchDuration:         time.Duration(environment.GetInt64(maxRelayBatchDuration, defaultBatchDuration)) * time.Second,
		maxServiceRecordBatchSize:     int(environment.GetInt64(maxServiceRecordBatchSize, defaultBatchSize)),
		maxServiceRecordBatchDuration: time.Duration(environment.GetInt64(maxServiceRecordBatchDuration, defaultBatchDuration)) * time.Second,
		dbTimeout:                     time.Duration(environment.GetInt64(dbTimeout, defaultDBTimeout)) * time.Second,
		debug:                         environment.GetBool(debug, defaultDebug),
		chanSize:                      int(environment.GetInt64(chanSize, defaultChanSize)),
	}
}

// CloudSQLConfig.GetDriver connects to a GCP CloudSQL instance using the cloudsqlconn lib.
// Intended for production use. Will be used if APP_ENV is 'production'.
func (c *cloudSQLConfig) GetDriver(ctx context.Context) (driver *postgresdriver.PostgresDriver, cleanup func() error, err error) {
	driverConfig := postgresdriver.CloudSQLConfig{
		DBUser:                 c.options.pgUser,
		DBPassword:             c.options.pgPassword,
		DBName:                 c.options.pgDatabase,
		InstanceConnectionName: c.options.dbInstanceConnectionName,
	}
	if c.options.privateIP {
		driverConfig.PrivateIP = "true"
	}

	driver, cleanup, err = postgresdriver.NewCloudSQLPostgresDriver(driverConfig)
	if err != nil {
		return nil, nil, err
	}

	return driver, cleanup, nil
}

// testDBConfig.GetDriver connects to a Postgres database using standard connection string and user/PW.
// Intended to be used for running tests on a local Docker container. Will be used if APP_ENV is 'test' or 'development'.
func (c *testDBConfig) GetDriver(ctx context.Context) (driver *postgresdriver.PostgresDriver, cleanup func() error, err error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.options.pgHost,
		c.options.pgPort,
		c.options.pgUser,
		c.options.pgPassword,
		c.options.pgDatabase,
	)

	driver, cleanup, err = postgresdriver.NewPostgresDriver(connectionString)
	if err != nil {
		return nil, nil, err
	}

	return driver, cleanup, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	options := gatherOptions()

	logConfig := zap.NewProductionConfig()
	logConfig.DisableStacktrace = true
	logConfig.DisableCaller = true
	logConfig.EncoderConfig.TimeKey = ""

	if options.debug {
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	log := zap.Must(logConfig.Build())

	// Choose DB configuration based on DB config vars
	var dbConfig DBConfig
	switch {
	// For CloudSQL DB
	case options.dbInstanceConnectionName != "":
		dbConfig = &cloudSQLConfig{options: options}

	// For local DB
	case options.pgHost != "" && options.pgPort != "":
		dbConfig = &testDBConfig{options: options}

	default:
		panic("invalid DB configuration")
	}

	driver, cleanup, err := dbConfig.GetDriver(context.Background())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Error(fmt.Sprintf("Failed to clean up: %v", err))
		}
	}()

	metricsExporter := metric.NewMetricExporter()

	relayBatch := batch.NewBatch(options.maxRelayBatchSize, options.chanSize, "relay", options.maxRelayBatchDuration, options.dbTimeout, driver.WriteRelays, log, metricsExporter)
	serviceRecordBatch := batch.NewBatch(options.maxServiceRecordBatchSize, options.chanSize, "service_record", options.maxServiceRecordBatchDuration, options.dbTimeout, driver.WriteServiceRecords, log, metricsExporter)

	router, err := router.NewRouter(driver, options.apiKeys, options.port, relayBatch, serviceRecordBatch, log)
	if err != nil {
		panic(err)
	}

	router.RunServer(ctx)
}
