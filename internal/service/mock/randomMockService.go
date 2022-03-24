package mock

import (
	"math/rand"
	"time"

	"github.com/ExchangeRates/produce-mock/internal/model"
)

type randomMockServiceImpl struct {
}

func NewRandomMockService() MockService {
	return &randomMockServiceImpl{}
}

func (m *randomMockServiceImpl) Generate(major, minor string, size int) []model.CupRatePoint {
	baseTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	points := make([]model.CupRatePoint, size)
	for i := 0; i < size; i++ {
		high := 1. + rand.Float64()
		low := high - rand.Float64()
		start := m.timeToMills(baseTime)
		end := m.timeToMills(baseTime.Add(time.Minute))
		points[i] = model.CupRatePoint{
			Major: major,
			Minor: minor,
			High:  high,
			Low:   low,
			Open:  high,
			Close: low,
			Start: start,
			End:   end,
		}
		baseTime = baseTime.Add(time.Second)
	}
	return points
}

func (m *randomMockServiceImpl) timeToMills(date time.Time) int64 {
	return int64(time.Nanosecond) * date.UnixNano() / int64(time.Microsecond)
}
