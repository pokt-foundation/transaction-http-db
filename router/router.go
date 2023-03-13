package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pokt-foundation/transaction-db/types"
	jsonresponse "github.com/pokt-foundation/utils-go/json-response"

	"github.com/sirupsen/logrus"
)

type driver interface {
	WriteSession(ctx context.Context, session types.PocketSession) error
	WriteRegion(ctx context.Context, region types.PortalRegion) error
	WriteRelay(ctx context.Context, relay types.Relay) error
	ReadRelay(ctx context.Context, relayID int) (types.Relay, error)
}

type Router struct {
	Router  *mux.Router
	Driver  driver
	APIKeys map[string]bool
	log     *logrus.Logger
}

func (rt *Router) logError(err error) {
	fields := logrus.Fields{
		"err": err.Error(),
	}

	rt.log.WithFields(fields).Error(err)
}

func respondWithResultOK(w http.ResponseWriter) {
	jsonresponse.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

// NewRouter returns router instance
func NewRouter(driver driver, apiKeys map[string]bool, logger *logrus.Logger) (*Router, error) {
	rt := &Router{
		Driver:  driver,
		Router:  mux.NewRouter(),
		APIKeys: apiKeys,
		log:     logger,
	}

	rt.Router.HandleFunc("/", rt.HealthCheck).Methods(http.MethodGet)

	rt.Router.HandleFunc("/v0/session", rt.CreateSession).Methods(http.MethodPost)
	rt.Router.HandleFunc("/v0/region", rt.CreateRegion).Methods(http.MethodPost)
	rt.Router.HandleFunc("/v0/relay", rt.CreateRelay).Methods(http.MethodPost)
	rt.Router.HandleFunc("/v0/relay/{id}", rt.GetRelay).Methods(http.MethodGet)

	rt.Router.Use(rt.AuthorizationHandler)

	return rt, nil
}

func (rt *Router) AuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is the path of the health check endpoint
		if r.URL.Path == "/" {
			h.ServeHTTP(w, r)

			return
		}

		if !rt.APIKeys[r.Header.Get("Authorization")] {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized"))
			if err != nil {
				panic(err)
			}

			return
		}

		h.ServeHTTP(w, r)
	})
}

func (rt *Router) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Transaction HTTP DB is up and running!"))
	if err != nil {
		panic(err)
	}
}

func (rt *Router) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	var session types.PocketSession
	err := decoder.Decode(&session)
	if err != nil {
		rt.logError(fmt.Errorf("CreateSession in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.Driver.WriteSession(ctx, session)
	if err != nil {
		rt.logError(fmt.Errorf("CreateSession in WriteSession failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) CreateRegion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	var region types.PortalRegion
	err := decoder.Decode(&region)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRegion in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.Driver.WriteRegion(ctx, region)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRegion in WriteRegion failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) CreateRelay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	var relay types.Relay
	err := decoder.Decode(&relay)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRelay in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.Driver.WriteRelay(ctx, relay)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRelay in WriteRelay failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) GetRelay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rt.logError(fmt.Errorf("GetRelay in params parsing failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	relay, err := rt.Driver.ReadRelay(ctx, id)
	if err != nil {
		rt.logError(fmt.Errorf("GetRelay in ReadRelay failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonresponse.RespondWithJSON(w, http.StatusOK, relay)
}
