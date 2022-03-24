package mock

import (
	"github.com/ExchangeRates/produce-mock/internal/model"
)

type MockService interface {
	Generate(major, minor string, size int) []model.CupRatePoint
}
