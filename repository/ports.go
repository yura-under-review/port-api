package repository

import (
	"context"
	"fmt"

	"github.com/yura-under-review/port-api/models"
	"github.com/yura-under-review/ports-domain-service/api"
	"google.golang.org/grpc"
)

type PortsRepository struct {
	address string
	client  api.PortsDomainServiceClient
	conn    *grpc.ClientConn
}

func New(config Config) *PortsRepository {
	return &PortsRepository{
		address: config.Address,
	}
}

func (r *PortsRepository) Init() error {
	var err error
	r.conn, err = grpc.Dial(r.address, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	r.client = api.NewPortsDomainServiceClient(r.conn)

	return nil
}

func (r *PortsRepository) Close() error {
	return r.conn.Close()
}

func (r *PortsRepository) UpsertPorts(ctx context.Context, ports []*models.PortInfo) error {

	apiPorts := ToAPIPorts(ports)

	_, err := r.client.BatchUpsertPorts(ctx, &api.BatchUpsertPortsRequest{
		Ports: apiPorts,
	})

	if err != nil {
		return fmt.Errorf("failed to upsert ports: %v", err)
	}

	return nil
}
