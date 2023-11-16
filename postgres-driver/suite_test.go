package postgresdriver

import (
	"context"
	"testing"
	"time"

	"github.com/pokt-foundation/transaction-http-db/types"

	"github.com/stretchr/testify/suite"
)

const (
	connectionString = "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable" // pragma: allowlist secret
)

type PGDriverTestSuite struct {
	suite.Suite
	connectionString string
	driver           *PostgresDriver

	// Records inserted on setup for testing purposes
	firstRelay         types.Relay
	firstServiceRecord types.ServiceRecord
}

func Test_RunPGDriverSuite(t *testing.T) {
	testSuite := new(PGDriverTestSuite)
	testSuite.connectionString = connectionString

	suite.Run(t, testSuite)
}

// SetupSuite runs before each test suite run
func (ts *PGDriverTestSuite) SetupSuite() {
	ts.NoError(ts.initPostgresDriver())

	ts.NoError(ts.driver.WriteRegion(context.Background(), types.PortalRegion{
		PortalRegionName: "La Colombia",
	}))

	ts.NoError(ts.driver.WriteSession(context.Background(), types.PocketSession{
		SessionKey:       "22",
		SessionHeight:    22,
		PortalRegionName: "La Colombia",
	}))

	ts.NoError(ts.driver.WriteRelay(context.Background(), types.Relay{
		PoktChainID:              "21",
		EndpointID:               "21",
		SessionKey:               "22",
		ProtocolAppPublicKey:     "21",
		RelaySourceURL:           "pablo.com",
		PoktNodeAddress:          "21",
		PoktNodeDomain:           "pablos.com",
		PoktNodePublicKey:        "aaa",
		RelayStartDatetime:       time.Now(),
		RelayReturnDatetime:      time.Now(),
		RelayRoundtripTime:       1,
		RelayChainMethodIDs:      []string{"get_height"},
		RelayDataSize:            21,
		RelayPortalTripTime:      21,
		RelayNodeTripTime:        21,
		RelayURLIsPublicEndpoint: false,
		PortalRegionName:         "La Colombia",
		IsAltruistRelay:          false,
		IsUserRelay:              false,
		RequestID:                "21",
	}))

	ts.NoError(ts.driver.WriteServiceRecord(context.Background(), types.ServiceRecord{
		NodePublicKey:          "21",
		PoktChainID:            "21",
		SessionKey:             "22",
		RequestID:              "21",
		PortalRegionName:       "La Colombia",
		Latency:                21.07,
		Tickets:                2,
		Result:                 "a",
		Available:              true,
		Successes:              21,
		Failures:               7,
		P90SuccessLatency:      21.07,
		MedianSuccessLatency:   21.07,
		WeightedSuccessLatency: 21.07,
		SuccessRate:            21,
	}))

	firstRelay, err := ts.driver.ReadRelay(context.Background(), 1)
	ts.NoError(err)

	firstServiceRecord, err := ts.driver.ReadServiceRecord(context.Background(), 1)
	ts.NoError(err)

	ts.firstRelay = firstRelay
	ts.firstServiceRecord = firstServiceRecord
}

// Initializes a real instance of the Postgres driver that connects to the test Postgres Docker container
func (ts *PGDriverTestSuite) initPostgresDriver() error {
	driver, _, err := NewPostgresDriver(ts.connectionString)
	if err != nil {
		return err
	}

	ts.driver = driver

	return nil
}
