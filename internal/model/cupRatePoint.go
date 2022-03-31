package model

type CupRatePoint struct {
	Major string  `json:"major"`
	Minor string  `json:"minor"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
	Start int64   `json:"start"`
	End   int64   `json:"end"`
}
