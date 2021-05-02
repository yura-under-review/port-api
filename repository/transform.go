package repository

import (
	"github.com/yura-under-review/port-api/models"
	"github.com/yura-under-review/ports-domain-service/api"
)

func ToAPIPort(src models.Port) *api.PortInfo {

	// TODO: implement
	return &api.PortInfo{}
}

func ToAPIPorts(src []models.Port) []*api.PortInfo {

	res := make([]*api.PortInfo, 0, len(src))

	for _, port := range src {
		res = append(res, ToAPIPort(port))
	}

	return res
}
