package transport

import (
	"context"
	pb "github.com/wahyunurdian26/product-service/contract/client"
)

func decodeListProductsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq.(*pb.ListProductsRequest), nil
}

func encodeListProductsResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pb.ListProductsResponse), nil
}
