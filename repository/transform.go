package repository

import (
	"github.com/yura-under-review/port-api/models"
	"github.com/yura-under-review/ports-domain-service/api"
)

func ToAPIPort(src *models.PortInfo) *api.PortInfo {

	var coord *api.Coordinate

	if len(src.Coordinate) >= 1 {
		coord = &api.Coordinate{
			Latitude:  src.Coordinate[0],
			Longitude: src.Coordinate[1],
		}
	}
	return &api.PortInfo{
		Symbol:     src.Symbol,
		Name:       src.Name,
		City:       src.City,
		Province:   src.Province,
		Country:    src.Country,
		Alias:      src.Alias,
		Regions:    src.Regions,
		Timezones:  src.Timezones,
		Unlocks:    src.Unlocks,
		Code:       src.Code,
		Coordinate: coord,
	}
}

func ToAPIPorts(src []*models.PortInfo) []*api.PortInfo {

	res := make([]*api.PortInfo, 0, len(src))

	for _, port := range src {
		res = append(res, ToAPIPort(port))
	}

	return res
}
