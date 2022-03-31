package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ExchangeRates/produce-mock/internal/kafka"
	"github.com/ExchangeRates/produce-mock/internal/service/mock"
	"github.com/sirupsen/logrus"
)

type MockController struct {
	log         *logrus.Logger
	service     mock.MockService
	cupProducer *kafka.ProducerCupRatePoint
}

func NewMockController(service mock.MockService, cupProducer *kafka.ProducerCupRatePoint) *MockController {
	return &MockController{
		log:         logrus.New(),
		service:     service,
		cupProducer: cupProducer,
	}
}

func (c *MockController) HandleRandomMock() http.HandlerFunc {
	type request struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
		Size  int    `json:"size"`
	}
	type response struct{}

	return func(w http.ResponseWriter, r *http.Request) {
		body := &request{}
		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			// TODO add processing of error
			c.log.Error(err)
			return
		}

		points := c.service.Generate(body.Major, body.Minor, body.Size)
		c.cupProducer.ProduceAll(points)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{}); err != nil {
			// TODO add processing of error
			c.log.Error(err)
			return
		}
	}
}
