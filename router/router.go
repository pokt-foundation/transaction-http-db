package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/pokt-foundation/transaction-db/types"
	"github.com/pokt-foundation/transaction-http-db/batch"
	jsonresponse "github.com/pokt-foundation/utils-go/json-response"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Driver interface {
	WriteSession(ctx context.Context, session types.PocketSession) error
	WriteRegion(ctx context.Context, region types.PortalRegion) error
	WriteRelay(ctx context.Context, relay types.Relay) error
	ReadRelay(ctx context.Context, relayID int) (types.Relay, error)
	WriteServiceRecord(ctx context.Context, serviceRecord types.ServiceRecord) error
	ReadServiceRecord(ctx context.Context, serviceRecordID int) (types.ServiceRecord, error)
}

type Router struct {
	router             *mux.Router
	driver             Driver
	apiKeys            map[string]bool
	relayBatch         *batch.Batch[*types.Relay]
	serviceRecordBatch *batch.Batch[*types.ServiceRecord]
	port               string
	log                *zap.Logger
}

func (rt *Router) logError(err error) {
	rt.log.Error(err.Error(), zap.String("err", err.Error()))
}

func respondWithResultOK(w http.ResponseWriter) {
	jsonresponse.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

// NewRouter returns router instance
func NewRouter(driver Driver, apiKeys map[string]bool, port string, relayBatch *batch.Batch[*types.Relay], serviceRecordBatch *batch.Batch[*types.ServiceRecord], logger *zap.Logger) (*Router, error) {
	rt := &Router{
		driver:             driver,
		router:             mux.NewRouter(),
		apiKeys:            apiKeys,
		relayBatch:         relayBatch,
		serviceRecordBatch: serviceRecordBatch,
		port:               port,
		log:                logger,
	}

	rt.router.HandleFunc("/", rt.HealthCheck).Methods(http.MethodGet)

	rt.router.HandleFunc("/v0/session", rt.CreateSession).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/region", rt.CreateRegion).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/relay", rt.CreateRelay).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/relays", rt.CreateRelays).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/relay/{id}", rt.GetRelay).Methods(http.MethodGet)
	rt.router.HandleFunc("/v0/service-record", rt.CreateServiceRecord).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/service-records", rt.CreateServiceRecords).Methods(http.MethodPost)
	rt.router.HandleFunc("/v0/service-record/{id}", rt.GetServiceRecord).Methods(http.MethodGet)

	rt.router.Use(rt.AuthorizationHandler)

	return rt, nil
}

func (rt *Router) RunServer(ctx context.Context) {
	httpServer := &http.Server{
		Addr:    ":" + rt.port,
		Handler: rt.router,
	}

	rt.log.Info(fmt.Sprintf("Transaction HTTP DB running in port: %s", rt.port))

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		rt.log.Info("HTTP router context finished")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			rt.logError(fmt.Errorf("Error closing http server: %s", err))
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := rt.relayBatch.Save(); err != nil {
				rt.logError(fmt.Errorf("Error saving relay batch: %s", err))
			}
		}()
		go func() {
			defer wg.Done()
			if err := rt.serviceRecordBatch.Save(); err != nil {
				rt.logError(fmt.Errorf("Error saving service record batch: %s", err))
			}
		}()
		wg.Wait()

		return nil
	})

	if err := g.Wait(); err != nil {
		rt.log.Info(fmt.Sprintf("exit reason: %s", err.Error()))
	}
}

func (rt *Router) AuthorizationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is the path of the health check endpoint
		if r.URL.Path == "/" {
			h.ServeHTTP(w, r)

			return
		}

		if !rt.apiKeys[r.Header.Get("Authorization")] {
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

	if err := session.Validate(); err != nil {
		rt.logError(fmt.Errorf("CreateSession in validate session failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = rt.driver.WriteSession(ctx, session)
	if err != nil {
		rt.logError(fmt.Errorf("CreateSession in WriteSession failed: %w", err))

		if errors.Is(err, types.ErrRepeatedSessionKey) {
			jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

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

	err = rt.driver.WriteRegion(ctx, region)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRegion in WriteRegion failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) CreateRelay(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var relay types.Relay
	err := decoder.Decode(&relay)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRelay in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.relayBatch.Add(&relay)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRelay in relay validating failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) CreateRelays(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var relays []types.Relay
	err := decoder.Decode(&relays)
	if err != nil {
		rt.logError(fmt.Errorf("CreateRelays in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	errs := 0
	for _, relay := range relays {
		err = rt.relayBatch.Add(&relay)
		if err != nil {
			rt.logError(fmt.Errorf("CreateRelays in relay validating failed: %w", err))
			errs++
		}
	}

	// TODO: Return the relay errors that failed
	if errs > 0 {
		msg := fmt.Sprintf("not all relays were processed successfully. failed relays: %d", errs)
		jsonresponse.RespondWithError(w, http.StatusBadRequest, msg)
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

	relay, err := rt.driver.ReadRelay(ctx, id)
	if err != nil {
		rt.logError(fmt.Errorf("GetRelay in ReadRelay failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonresponse.RespondWithJSON(w, http.StatusOK, relay)
}

func (rt *Router) CreateServiceRecord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var serviceRecord types.ServiceRecord
	err := decoder.Decode(&serviceRecord)
	if err != nil {
		rt.logError(fmt.Errorf("CreateServiceRecord in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	err = rt.serviceRecordBatch.Add(&serviceRecord)
	if err != nil {
		rt.logError(fmt.Errorf("CreateServiceRecord in service record validating failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) CreateServiceRecords(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var serviceRecords []types.ServiceRecord
	err := decoder.Decode(&serviceRecords)
	if err != nil {
		rt.logError(fmt.Errorf("CreateServiceRecords in JSON decoding failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	errs := 0
	for _, serviceRecord := range serviceRecords {
		err = rt.serviceRecordBatch.Add(&serviceRecord)
		if err != nil {
			rt.logError(fmt.Errorf("CreateServiceRecords in service record validating failed: %w", err))
			errs++
		}
	}

	// TODO: Return the service records errors that failed
	if errs > 0 {
		msg := fmt.Sprintf("not all service records were processed successfully. failed service records: %d", errs)
		jsonresponse.RespondWithError(w, http.StatusBadRequest, msg)
		return
	}

	respondWithResultOK(w)
}

func (rt *Router) GetServiceRecord(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rt.logError(fmt.Errorf("GetServiceRecord in params parsing failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	serviceRecord, err := rt.driver.ReadServiceRecord(ctx, id)
	if err != nil {
		rt.logError(fmt.Errorf("GetServiceRecord in ReadServiceRecord failed: %w", err))
		jsonresponse.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonresponse.RespondWithJSON(w, http.StatusOK, serviceRecord)
}
