package main

import (
	"net/http"

	postgresdriver "github.com/pokt-foundation/transaction-db/postgres-driver"
	"github.com/pokt-foundation/transaction-http-db/router"
	"github.com/pokt-foundation/utils-go/environment"
	"github.com/sirupsen/logrus"
)

const (
	connectionString = "CONNECTION_STRING"
	apiKeys          = "API_KEYS"
	port             = "PORT"

	defaultPort = "8080"
)

type options struct {
	connectionString string
	apiKeys          map[string]bool
	port             string
}

func gatherOptions() options {
	return options{
		connectionString: environment.MustGetString(connectionString),
		apiKeys:          environment.MustGetStringMap(apiKeys, ","),
		port:             environment.GetString(port, defaultPort),
	}
}

func httpHandler(router *router.Router, port string, log *logrus.Logger) {
	log.Printf("Transaction HTTP DB running in port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router.Router))
}

func main() {
	log := logrus.New()
	// log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	options := gatherOptions()

	driver, err := postgresdriver.NewPostgresDriver(options.connectionString)
	if err != nil {
		panic(err)
	}

	router, err := router.NewRouter(driver, options.apiKeys, log)
	if err != nil {
		panic(err)
	}

	httpHandler(router, options.port, log)
}
