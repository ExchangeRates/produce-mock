package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ExchangeRates/produce-mock/internal/controller"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "gopherschool"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router         *mux.Router
	log            *logrus.Logger
	mockController *controller.MockController
}

func NewServer(mockController *controller.MockController) *server {
	s := &server{
		router:         mux.NewRouter(),
		log:            logrus.New(),
		mockController: mockController,
	}

	s.configureRouter()

	s.log.Info("Starting api server")

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	s.router.Path("/mock/random").
		Handler(s.mockController.HandleRandomMock()).
		Methods(http.MethodPost)
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.log.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("Started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		next.ServeHTTP(w, r)

		s.log.Logf(
			logrus.InfoLevel,
			"Completed with %v",
			time.Now().Sub(start),
		)

	})
}

func (s *server) BindingAddressFromPort(port int) string {
	s.log.WithFields(logrus.Fields{
		"port": port,
	}).Info("listening by port")
	return fmt.Sprintf(":%d", port)
}
