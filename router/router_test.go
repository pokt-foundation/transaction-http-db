package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pokt-foundation/transaction-db/types"

	"github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouter_HealthCheck(t *testing.T) {
	c := require.New(t)

	router, err := NewRouter(&mockDriver{}, map[string]bool{"": true}, logrus.New())
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
		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateError(t *testing.T) {
	c := require.New(t)

	driverMock := &mockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, logrus.New())
	c.NoError(err)

	rawErrorToSend := types.Error{
		ErrorCode:        404,
		ErrorName:        "not found",
		ErrorDescription: "we did not find it",
	}

	errToSend, err := json.Marshal(rawErrorToSend)
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
			reqInput:           errToSend,
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
			reqInput:            errToSend,
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
		req, err := http.NewRequest(http.MethodPost, "/v0/error", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("WriteError", mock.Anything, mock.Anything).Return(tt.errReturnedByDriver).Once()
		}

		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateSession(t *testing.T) {
	c := require.New(t)

	driverMock := &mockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, logrus.New())
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

		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRegion(t *testing.T) {
	c := require.New(t)

	driverMock := &mockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, logrus.New())
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

		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_CreateRelay(t *testing.T) {
	c := require.New(t)

	driverMock := &mockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, logrus.New())
	c.NoError(err)

	rawRelayToSend := types.Relay{
		RelayRoundtripTime: 21,
	}

	relayToSend, err := json.Marshal(rawRelayToSend)
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
			reqInput:           relayToSend,
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
			reqInput:            relayToSend,
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
		req, err := http.NewRequest(http.MethodPost, "/v0/relay", bytes.NewBuffer(tt.reqInput))
		c.NoError(err)

		req.Header.Set("Authorization", tt.apiKey)
		rr := httptest.NewRecorder()

		if tt.setMock {
			driverMock.On("WriteRelay", mock.Anything, mock.Anything).Return(tt.errReturnedByDriver).Once()
		}

		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
	}
}

func TestRouter_GetRelay(t *testing.T) {
	c := require.New(t)

	driverMock := &mockDriver{}
	router, err := NewRouter(driverMock, map[string]bool{"": true}, logrus.New())
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

		router.Router.ServeHTTP(rr, req)
		c.Equal(tt.expectedStatusCode, rr.Code)
		c.Equal(tt.expectedBody, rr.Body.String())
	}
}
