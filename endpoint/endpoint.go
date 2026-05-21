package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/wahyunurdian26/product-service/service"
	pb "github.com/wahyunurdian26/product-service/contract/client"
)

type Endpoints struct {
	CreateProductEndpoint endpoint.Endpoint
	ListProductsEndpoint  endpoint.Endpoint
	HealthProduct         endpoint.Endpoint
}

func MakeEndpoints(svc service.ProductService) Endpoints {
	return Endpoints{
		CreateProductEndpoint: makeCreateProductEndpoint(svc),
		ListProductsEndpoint:  makeListProductsEndpoint(svc),
		HealthProduct:         makeHealthProductEndpoint(svc),
	}
}

func makeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.CreateProductRequest)
		return svc.CreateProduct(ctx, req)
	}
}

func makeListProductsEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.ListProductsRequest)
		return svc.ListProducts(ctx, req)
	}
}

func makeHealthProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*pb.HealthProductRequest)
		if !ok {
			req = &pb.HealthProductRequest{}
		}
		return svc.HealthProduct(ctx, req)
	}
}
