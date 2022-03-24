package api

import (
	"fmt"
	"net/http"

	"github.com/ExchangeRates/produce-mock/internal/controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router         *mux.Router
	logger         *logrus.Logger
	mockController *controller.MockController
}

func NewServer(mockController *controller.MockController) *server {
	s := &server{
		router:         mux.NewRouter(),
		logger:         logrus.New(),
		mockController: mockController,
	}

	s.configureRouter()

	logrus.Info("Starting api server")

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Path("/mock/random").
		Handler(s.mockController.HandleRandomMock()).
		Methods(http.MethodPost)
}

func (s *server) BindingAddressFromPort(port int) string {
	logrus.Info("listening on port ", port)
	return fmt.Sprintf(":%d", port)
}
