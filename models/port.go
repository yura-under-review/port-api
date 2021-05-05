package models

type (
	PortInfo struct {
		Symbol     string
		Name       string    `json:"name"`
		City       string    `json:"city"`
		Province   string    `json:"province"`
		Country    string    `json:"country"`
		Alias      []string  `json:"alias"`
		Regions    []string  `json:"regions"`
		Timezones  []string  `json:"timezones"`
		Unlocks    []string  `json:"unlocs"`
		Code       string    `json:"code"`
		Coordinate []float32 `json:"coordinates"`
	}
)
