package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/wahyunurdian26/product-service/model"
	"github.com/wahyunurdian26/product-service/service"
	"github.com/wahyunurdian26/product-service/validation"
)

type Endpoints struct {
	CreateProduct endpoint.Endpoint
	ListProducts  endpoint.Endpoint
}

func MakeEndpoints(svc service.ProductService) Endpoints {
	return Endpoints{
		CreateProduct: makeCreateProductEndpoint(svc),
		ListProducts:  makeListProductsEndpoint(svc),
	}
}

func makeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*model.CreateProductRequest)
		
		if err := validation.CreateProductValidation(req); err != nil {
			return model.CreateProductResponse{Error: err.Error()}, nil
		}

		product, err := svc.CreateProduct(ctx, req)
		if err != nil {
			return model.CreateProductResponse{Error: err.Error()}, nil
		}
		return model.CreateProductResponse{Product: product, Message: "Success"}, nil
	}
}

func makeListProductsEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*model.ListProductRequest)
		products, err := svc.ListProducts(ctx, *req)
		if err != nil {
			return model.ListProductResponse{Error: err.Error()}, nil
		}
		return model.ListProductResponse{Products: products, Message: "Success"}, nil
	}
}
