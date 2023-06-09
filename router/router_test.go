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
	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
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

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
	c.NoError(err)

	rawSessionToSend := types.PocketSession{
		SessionKey:       "21",
		SessionHeight:    1,
		PortalRegionName: "region",
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
		repeat              bool
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			reqInput:           sessionToSend,
			setMock:            true,
			repeat:             true,
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

		if tt.repeat {
			req, err := http.NewRequest(http.MethodPost, "/v0/session", bytes.NewBuffer(tt.reqInput))
			c.NoError(err)

			req.Header.Set("Authorization", tt.apiKey)
			rr := httptest.NewRecorder()

			driverMock.On("WriteSession", mock.Anything, mock.Anything).Return(types.ErrRepeatedSessionKey).Once()

			router.router.ServeHTTP(rr, req)
			c.Equal(http.StatusBadRequest, rr.Code)
		}
	}
}

func TestRouter_CreateRegion(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
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

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())
	relayWriterMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
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

func TestRouter_CreateServiceRecord(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())
	serviceRecordMock.On("WriteServiceRecords", mock.Anything, mock.Anything).Return(nil).Once()

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
	c.NoError(err)

	rawServiceRecordToSend := types.ServiceRecord{
		SessionKey:             "21",
		NodePublicKey:          "21",
		PoktChainID:            "21",
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
	}

	serviceRecordToSend, err := json.Marshal(rawServiceRecordToSend)
	c.NoError(err)

	rawWrongServiceRecordToSend := types.ServiceRecord{
		SessionKey: "1",
	}

	wrongServiceRecordToSend, err := json.Marshal(rawWrongServiceRecordToSend)
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
			reqInput:           serviceRecordToSend,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:               "Invalid Relay",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           wrongServiceRecordToSend,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/service-record", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRelays(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())
	relayWriterMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
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

	rawWrongRelayToSend := []types.Relay{{
		PoktChainID: "21",
	}}

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

func TestRouter_CreateServiceRecords(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())
	serviceRecordMock.On("WriteServiceRecords", mock.Anything, mock.Anything).Return(nil).Once()

	router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
	c.NoError(err)

	rawServiceRecordsToSend := []types.ServiceRecord{{
		SessionKey:             "21",
		NodePublicKey:          "21",
		PoktChainID:            "21",
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
	},
		{
			SessionKey:             "21",
			NodePublicKey:          "21",
			PoktChainID:            "21",
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
		}}

	serviceRecordToSend, err := json.Marshal(rawServiceRecordsToSend)
	c.NoError(err)

	rawWrongServiceRecordToSend := []types.ServiceRecord{{
		SessionKey: "1",
	}}

	wrongServiceRecordToSend, err := json.Marshal(rawWrongServiceRecordToSend)
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
			reqInput:           serviceRecordToSend,
		},
		{
			name:               "Wrong input",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           []byte("wrong"),
		},
		{
			name:               "Invalid Relay",
			expectedStatusCode: http.StatusBadRequest,
			reqInput:           wrongServiceRecordToSend,
		},
		{
			name:               "Not authorized",
			expectedStatusCode: http.StatusUnauthorized,
			apiKey:             "wrong",
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(http.MethodPost, "/v0/service-records", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_GetRelay(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
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

func TestRouter_GetServiceRecord(t *testing.T) {
	c := require.New(t)

	relayWriterMock := &batch.MockRelayWriter{}
	relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

	serviceRecordMock := &batch.MockServiceRecordWriter{}
	serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

	driverMock := &MockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
	c.NoError(err)

	serviceRecordToReturn := types.ServiceRecord{
		SessionKey: "1",
	}

	expectedBody, err := json.Marshal(serviceRecordToReturn)
	c.NoError(err)

	bodyString := string(expectedBody)

	tests := []struct {
		name                          string
		reqInput                      string
		expectedStatusCode            int
		expectedBody                  string
		serviceRecordReturnedByDriver types.ServiceRecord
		errReturnedByDriver           error
		apiKey                        string
		setMock                       bool
	}{
		{
			name:                          "Success",
			reqInput:                      "1",
			expectedStatusCode:            http.StatusOK,
			expectedBody:                  bodyString,
			serviceRecordReturnedByDriver: serviceRecordToReturn,
			setMock:                       true,
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
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v0/service-record/%s", tt.reqInput), nil)
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("ReadServiceRecord", mock.Anything, mock.Anything).Return(tt.serviceRecordReturnedByDriver, tt.errReturnedByDriver).Once()
		}

		router.router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
		c.Equal(tt.expectedBody, rr.Body.String())
	}
}

func TestRouter_RunServer(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name                      string
		ctxTimeout                time.Duration
		expectedRelaysSize        int
		expectedServicRecordsSize int
	}{
		{
			name:                      "Context finished",
			ctxTimeout:                time.Millisecond,
			expectedRelaysSize:        0,
			expectedServicRecordsSize: 0,
		},
		{
			name:                      "Context not finished",
			ctxTimeout:                time.Minute,
			expectedRelaysSize:        1,
			expectedServicRecordsSize: 1,
		},
	}

	for _, tt := range tests {
		relayWriterMock := &batch.MockRelayWriter{}
		relayBatch := batch.NewBatch(2, "relay", time.Hour, time.Hour, relayWriterMock.WriteRelays, logrus.New())

		serviceRecordMock := &batch.MockServiceRecordWriter{}
		serviceRecordBatch := batch.NewBatch(2, "service_record", time.Hour, time.Hour, serviceRecordMock.WriteServiceRecords, logrus.New())

		err := relayBatch.Add(&types.Relay{
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

		err = serviceRecordBatch.Add(&types.ServiceRecord{
			SessionKey:             "21",
			NodePublicKey:          "21",
			PoktChainID:            "21",
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
		})
		c.NoError(err)

		time.Sleep(time.Second)
		c.Equal(1, relayBatch.Size())
		c.Equal(1, serviceRecordBatch.Size())

		router, err := NewRouter(&MockDriver{}, map[string]bool{"": true}, "8080", relayBatch, serviceRecordBatch, logrus.New())
		c.NoError(err)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), tt.ctxTimeout)
		defer cancel()

		relayWriterMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()
		serviceRecordMock.On("WriteServiceRecords", mock.Anything, mock.Anything).Return(nil).Once()

		go router.RunServer(ctxTimeout)

		time.Sleep(time.Second)
		c.Equal(tt.expectedRelaysSize, relayBatch.Size())
		c.Equal(tt.expectedServicRecordsSize, serviceRecordBatch.Size())
	}
}
