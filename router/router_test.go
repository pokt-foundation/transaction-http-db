package router

import (
	"bytes"
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pokt-foundation/transaction-db/types"
	"github.com/pokt-foundation/transaction-http-db/batch"
	"github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouter_HealthCheck(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(2, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		c.NoError(err)

		rr := httptest.NewRecorder()
		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateSession(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(2, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	rawSessionToSend := types.PocketSession{
		SessionKey: "21",
	}

	sessionToSend, err := json.Marshal(rawSessionToSend)
	c.NoError(err)

	tests := []struct {
		name                string
		expectedStatusCode  int
		reqInput            []byte
		errReturnedByDriver error
		apiKey              string
		setMock             bool
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           sessionToSend,
			setMock:            true,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:                "Failure on driver",
			expectedStatusCode:  http.StatusInternalServerError,
			reqInput:            sessionToSend,
			errReturnedByDriver: errors.New("dummy"),
			setMock:             true,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/session", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("WriteSession", mock.Anything, mock.Anything).Return(tt.errReturnedByDriver).Once()
		}

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRegion(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(2, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	rawRegionToSend := types.PortalRegion{
		PortalRegionName: "Los Praditos",
	}

	regionToSend, err := json.Marshal(rawRegionToSend)
	c.NoError(err)

	tests := []struct {
		name                string
		expectedStatusCode  int
		reqInput            []byte
		errReturnedByDriver error
		apiKey              string
		setMock             bool
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           regionToSend,
			setMock:            true,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:                "Failure on driver",
			expectedStatusCode:  http.StatusInternalServerError,
			reqInput:            regionToSend,
			errReturnedByDriver: errors.New("dummy"),
			setMock:             true,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/region", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("WriteRegion", mock.Anything, mock.Anything).Return(tt.errReturnedByDriver).Once()
		}

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRelay(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(2, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	rawRelayToSend := types.Relay{
		PoktChainID:              "21",
		EndpointID:               "21",
		SessionKey:               "21",
		ProtocolAppPublicKey:     "21",
		RelaySourceURL:           "pablo.com",
		PoktNodeAddress:          "21",
		PoktNodeDomain:           "pablos.com",
		PoktNodePublicKey:        "aaa",
		RelayStartDatetime:       time.Now(),
		RelayReturnDatetime:      time.Now(),
		IsError:                  true,
		ErrorCode:                21,
		ErrorName:                "favorite number",
		ErrorMessage:             "just Pablo can use it",
		ErrorType:                "chain_check",
		ErrorSource:              "internal",
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
		PoktTxID:                 "21",
	}

	relayToSend, err := json.Marshal(rawRelayToSend)
	c.NoError(err)

	rawWrongRelayToSend := types.Relay{
		PoktChainID: "21",
	}

	wrongRelayToSend, err := json.Marshal(rawWrongRelayToSend)
	c.NoError(err)

	tests := []struct {
		name               string
		expectedStatusCode int
		reqInput           []byte
		apiKey             string
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           relayToSend,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:               "Invalid Relay",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           wrongRelayToSend,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/relay", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRelays(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(3, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	rawRelaysToSend := []types.Relay{{
		PoktChainID:              "21",
		EndpointID:               "21",
		SessionKey:               "21",
		ProtocolAppPublicKey:     "21",
		RelaySourceURL:           "pablo.com",
		PoktNodeAddress:          "21",
		PoktNodeDomain:           "pablos.com",
		PoktNodePublicKey:        "aaa",
		RelayStartDatetime:       time.Now(),
		RelayReturnDatetime:      time.Now(),
		IsError:                  true,
		ErrorCode:                21,
		ErrorName:                "favorite number",
		ErrorMessage:             "just Pablo can use it",
		ErrorType:                "chain_check",
		ErrorSource:              "internal",
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
		PoktTxID:                 "21",
	},
		{
			PoktChainID:              "21",
			EndpointID:               "21",
			SessionKey:               "21",
			ProtocolAppPublicKey:     "21",
			RelaySourceURL:           "pablo.com",
			PoktNodeAddress:          "21",
			PoktNodeDomain:           "pablos.com",
			PoktNodePublicKey:        "aaa",
			RelayStartDatetime:       time.Now(),
			RelayReturnDatetime:      time.Now(),
			IsError:                  true,
			ErrorCode:                21,
			ErrorName:                "favorite number",
			ErrorMessage:             "just Pablo can use it",
			ErrorType:                "chain_check",
			ErrorSource:              "internal",
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
			PoktTxID:                 "21",
		}}

	relayToSend, err := json.Marshal(rawRelaysToSend)
	c.NoError(err)

	rawWrongRelayToSend := types.Relay{
		PoktChainID: "21",
	}

	wrongRelayToSend, err := json.Marshal(rawWrongRelayToSend)
	c.NoError(err)

	tests := []struct {
		name               string
		expectedStatusCode int
		reqInput           []byte
		apiKey             string
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           relayToSend,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:               "Invalid Relay",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           wrongRelayToSend,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/relays", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_GetRelay(t *testing.T) {
	c := require.New(t)

	batch := batch.NewRelayBatch(2, time.Hour, time.Hour, &batch.MockRelayWriter{}, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", batch, logrus.New())
	c.NoError(err)

	relayToReturn := types.Relay{
		RelayRoundtripTime: 21,
	}

	expectedBody, err := json.Marshal(relayToReturn)
	c.NoError(err)

	bodyString := string(expectedBody)

	tests := []struct {
		name                  string
		reqInput              string
		expectedStatusCode    int
		expectedBody          string
		relayReturnedByDriver types.Relay
		errReturnedByDriver   error
		apiKey                string
		setMock               bool
	}{
		{
			name:                  "Success",
			reqInput:              "1",
			expectedStatusCode:    http.StatusOK,
			expectedBody:          bodyString,
			relayReturnedByDriver: relayToReturn,
			setMock:               true,
		},
		{
			name:               "Wrong input",
			reqInput:           "pablo",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"strconv.Atoi: parsing \"pablo\": invalid syntax"}`,
		},
		{
			name:                "Failure on driver",
			reqInput:            "1",
			expectedStatusCode:  http.StatusInternalServerError,
			expectedBody:        `{"error":"dummy"}`,
			errReturnedByDriver: errors.New("dummy"),
			setMock:             true,
		},
		{
			name:               "Not authorized",
			reqInput:           "1",
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Unauthorized",
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v0/relay/%s", tt.reqInput), nil)
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("ReadRelay", mock.Anything, mock.Anything).Return(tt.relayReturnedByDriver, tt.errReturnedByDriver).Once()
		}

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
		c.Equal(tt.expectedBody, rr.Body.String())
	}
}

func TestRouter_RunServer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name               string
		ctxTimeout         time.Duration
		expectedRelaysSize int
	}{
		{
			name:               "Context finished",
			ctxTimeout:         time.Millisecond,
			expectedRelaysSize: 0,
		},
		{
			name:               "Context not finished",
			ctxTimeout:         time.Minute,
			expectedRelaysSize: 1,
		},
	}

	for _, tt := range tests {
		writerMock := &batch.MockRelayWriter{}
		batch := batch.NewRelayBatch(2, time.Hour, time.Hour, writerMock, logrus.New())

		err := batch.AddRelay(types.Relay{
			PoktChainID:              "21",
			EndpointID:               "21",
			SessionKey:               "21",
			ProtocolAppPublicKey:     "21",
			RelaySourceURL:           "pablo.com",
			PoktNodeAddress:          "21",
			PoktNodeDomain:           "pablos.com",
			PoktNodePublicKey:        "aaa",
			RelayStartDatetime:       time.Now(),
			RelayReturnDatetime:      time.Now(),
			IsError:                  true,
			ErrorCode:                21,
			ErrorName:                "favorite number",
			ErrorMessage:             "just Pablo can use it",
			ErrorType:                "chain_check",
			ErrorSource:              "internal",
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
			PoktTxID:                 "21",
		})
		c.NoError(err)

		time.Sleep(time.Second)
		c.Equal(1, batch.RelaysSize())

		router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", batch, logrus.New())
		c.NoError(err)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), tt.ctxTimeout)
		defer cancel()

		writerMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()

		go router.RunServer(ctxTimeout)

		time.Sleep(time.Second)
		c.Equal(tt.expectedRelaysSize, batch.RelaysSize())
	}
}
