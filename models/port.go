package models

type (
	PortInfo struct {
		Symbol     string
		Name       string
		City       string
		Province   string
		Country    string
		Alias      []string
		Regions    []string
		Timezones  []string
		Unlocks    []string
		Code       int
		Coordinate Coordinate
	}

	Coordinate struct {
		Latitude  float64
		Longitude float64
	}
)
