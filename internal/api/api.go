package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ExchangeRates/produce-mock/internal/controller"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	sessionName        = "gopherschool"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type Server interface {
	GracefullListenAndServe(port int) error
}

type server struct {
	router         *mux.Router
	log            *logrus.Logger
	mockController *controller.MockController
}

func NewServer(mockController *controller.MockController) Server {
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

func (s *server) bindingAddressFromPort(port int) string {
	s.log.WithFields(logrus.Fields{
		"port": port,
	}).Info("listening by port")
	return fmt.Sprintf(":%d", port)
}

func (s *server) GracefullListenAndServe(port int) error {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr:    s.bindingAddressFromPort(port),
		Handler: s,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	return g.Wait()
}
