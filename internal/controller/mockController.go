package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ExchangeRates/produce-mock/internal/service/mock"
	"github.com/sirupsen/logrus"
)

type MockController struct {
	service mock.MockService
}

func NewMockController(service mock.MockService) *MockController {
	return &MockController{
		service: service,
	}
}

func (m *MockController) HandleRandomMock() http.HandlerFunc {
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
			return
		}

		points := m.service.Generate(body.Major, body.Minor, body.Size)
		for _, point := range points {
			logrus.Info(fmt.Sprintf("%s:%s h: %f | l: %f (o: %f | c: %f) [%d - %d]",
				point.Major, point.Minor,
				point.High, point.Low,
				point.Open, point.Close,
				point.Start, point.End,
			))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{}); err != nil {
			// TODO add processing of error
			return
		}
	}
}
