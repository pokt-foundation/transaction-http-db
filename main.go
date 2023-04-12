package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	postgresdriver "github.com/pokt-foundation/transaction-db/postgres-driver"
	"github.com/pokt-foundation/transaction-http-db/batch"
	"github.com/pokt-foundation/transaction-http-db/router"
	"github.com/pokt-foundation/utils-go/environment"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	connectionString      = "CONNECTION_STRING"
	apiKeys               = "API_KEYS"
	port                  = "PORT"
	maxRelayBatchSize     = "MAX_RELAY_BATCH_SIZE"
	maxRelayBatchDuration = "MAX_RELAY_BATCH_DURATION"
	dbTimeout             = "DB_TIMEOUT"
	debug                 = "DEBUG"
	useSSH                = "USE_SSH"
	sshHost               = "SSH_HOST"
	sshPort               = "SSH_PORT"
	sshUser               = "SSH_USER"
	sshKeyFilePath        = "SSH_KEY_FILE_PATH"

	defaultPort          = "8080"
	defaultBatchSize     = 1000
	defaultBatchDuration = 60
	defaultDBTimeout     = 60
	defaultDebug         = false
	defaultUseSSH        = false
)

type options struct {
	connectionString      string
	apiKeys               map[string]bool
	port                  string
	maxRelayBatchSize     int
	maxRelayBatchDuration time.Duration
	dbTimeout             time.Duration
	debug                 bool
	useSSH                bool
	sshHost               string
	sshPort               string
	sshUser               string
	sshKeyFilePath        string
}

func gatherOptions() options {
	return options{
		connectionString:      environment.MustGetString(connectionString),
		apiKeys:               environment.MustGetStringMap(apiKeys, ","),
		port:                  environment.GetString(port, defaultPort),
		maxRelayBatchSize:     int(environment.GetInt64(maxRelayBatchSize, defaultBatchSize)),
		maxRelayBatchDuration: time.Duration(environment.GetInt64(maxRelayBatchDuration, defaultBatchDuration)) * time.Second,
		dbTimeout:             time.Duration(environment.GetInt64(dbTimeout, defaultDBTimeout)) * time.Second,
		debug:                 environment.GetBool(debug, defaultDebug),
		useSSH:                environment.GetBool(useSSH, defaultUseSSH),
		sshHost:               environment.GetString(sshHost, ""),
		sshPort:               environment.GetString(sshPort, ""),
		sshUser:               environment.GetString(sshUser, ""),
		sshKeyFilePath:        environment.GetString(sshKeyFilePath, ""),
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

	var driver *postgresdriver.PostgresDriver

	if options.useSSH {
		sshKey, err := ioutil.ReadFile(options.sshKeyFilePath)
		if err != nil {
			panic(err)
		}
		sshKeySigner, err := ssh.ParsePrivateKey(sshKey)
		if err != nil {
			panic(err)
		}

		sshConfig := &ssh.ClientConfig{
			User: options.sshUser,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(sshKeySigner),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		// Connect to the SSH Server
		sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", options.sshHost, options.sshPort), sshConfig)
		if err != nil {
			panic(err)
		}
		defer sshcon.Close()

		driver, err = postgresdriver.NewPostgresDriverWithSSH(options.connectionString, sshcon)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		driver, err = postgresdriver.NewPostgresDriver(options.connectionString)
		if err != nil {
			panic(err)
		}
	}

	relayBatch := batch.NewRelayBatch(options.maxRelayBatchSize, options.maxRelayBatchDuration, options.dbTimeout, driver, log)

	router, err := router.NewRouter(driver, options.apiKeys, options.port, relayBatch, log)
	if err != nil {
		panic(err)
	}

	router.RunServer(ctx)
}
