package internal

import (
	"net/http"

	"github.com/ExchangeRates/produce-mock/internal/api"
	"github.com/ExchangeRates/produce-mock/internal/config"
	"github.com/ExchangeRates/produce-mock/internal/controller"
	"github.com/ExchangeRates/produce-mock/internal/kafka"
	"github.com/ExchangeRates/produce-mock/internal/service/mock"
)

func Start(config *config.Config) error {

	mockService := mock.NewRandomMockService()
	cupRateProducer := kafka.NewProducerCupRatePoint(config.Urls)
	mockController := controller.NewMockController(mockService, cupRateProducer)

	srv := api.NewServer(mockController)
	bindingAddr := srv.BindingAddressFromPort(config.Port)

	return http.ListenAndServe(bindingAddr, srv)
}
