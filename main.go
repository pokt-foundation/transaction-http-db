package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgresdriver "github.com/pokt-foundation/transaction-db/postgres-driver"
	"github.com/pokt-foundation/transaction-http-db/batch"
	"github.com/pokt-foundation/transaction-http-db/router"
	"github.com/pokt-foundation/utils-go/environment"
	"github.com/sirupsen/logrus"
)

const (
	connectionString = "CONNECTION_STRING"
	apiKeys          = "API_KEYS"
	port             = "PORT"
	maxBatchSize     = "MAX_BATCH_SIZE"
	maxBatchDuration = "MAX_BATCH_DURATION"
	dbTimeout        = "DB_TIMEOUT"
	debug            = "DEBUG"

	defaultPort          = "8080"
	defaultBatchSize     = 1000
	defaultBatchDuration = 60
	defaultDBTimeout     = 60
	defaultDebug         = false
)

type options struct {
	connectionString string
	apiKeys          map[string]bool
	port             string
	maxBatchSize     int
	maxBatchDuration time.Duration
	dbTimeout        time.Duration
	debug            bool
}

func gatherOptions() options {
	return options{
		connectionString: environment.MustGetString(connectionString),
		apiKeys:          environment.MustGetStringMap(apiKeys, ","),
		port:             environment.GetString(port, defaultPort),
		maxBatchSize:     int(environment.GetInt64(maxBatchSize, defaultBatchSize)),
		maxBatchDuration: time.Duration(environment.GetInt64(maxBatchDuration, defaultBatchDuration)) * time.Second,
		dbTimeout:        time.Duration(environment.GetInt64(dbTimeout, defaultDBTimeout)) * time.Second,
		debug:            environment.GetBool(debug, defaultDebug),
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	options := gatherOptions()

	log := logrus.New()
	// log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	if options.debug {
		log.Level = logrus.DebugLevel
	}

	driver, err := postgresdriver.NewPostgresDriver(options.connectionString)
	if err != nil {
		panic(err)
	}

	batch := batch.New(options.maxBatchSize, options.maxBatchDuration, options.dbTimeout, driver, log)

	router, err := router.NewRouter(driver, options.apiKeys, options.port, batch, log)
	if err != nil {
		panic(err)
	}

	router.RunServer(ctx)
}
