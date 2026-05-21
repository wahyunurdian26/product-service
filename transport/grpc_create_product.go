package transport

import (
	"context"
	pb "github.com/wahyunurdian26/product-service/contract/client"
)

func decodeCreateProductRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.CreateProductRequest), nil
}

func encodeCreateProductResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pb.CreateProductResponse), nil
}

func decodeHealthCheckRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return nil, nil
}

func encodeHealthCheckResponse(_ context.Context, response interface{}) (interface{}, error) {
	return nil, nil // Not heavily used in simple implementation, rely on grpc_health_v1
}
