package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/Bifrost-Mesh/users-microservice/pkg/healthcheck"
)

type HealthcheckService struct {
	healthcheckables []healthcheck.Healthcheckable
}

func (h *HealthcheckService) Check(ctx context.Context,
	request *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	response := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
	}

	err := healthcheck.Healthcheck(h.healthcheckables)
	if err != nil {
		response.Status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
		return response, err
	}

	response.Status = grpc_health_v1.HealthCheckResponse_SERVING
	return response, nil
}

func (h *HealthcheckService) List(ctx context.Context,
	request *grpc_health_v1.HealthListRequest,
) (*grpc_health_v1.HealthListResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (h *HealthcheckService) Watch(
	request *grpc_health_v1.HealthCheckRequest,
	responseStream grpc.ServerStreamingServer[grpc_health_v1.HealthCheckResponse],
) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
