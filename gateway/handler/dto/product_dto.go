package dto

import (
	"github.com/wahyunurdian26/product-service/gateway/kit"
	"github.com/wahyunurdian26/product-service/model"
	"github.com/wahyunurdian26/product-service/validation"
)

type ParseCreateProduct struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Type  string  `json:"type"`
}

func ParseCreateProductRequest(ctx kit.Context) (ParseCreateProduct, error) {
	var req ParseCreateProduct
	if err := ctx.BindJSON(&req); err != nil {
		return req, err
	}
	return req, nil
}

func (req ParseCreateProduct) ToServiceRequest() (*model.CreateProductRequest, error) {
	serviceReq := &model.CreateProductRequest{
		Name:  req.Name,
		Price: req.Price,
		Type:  model.ProductType(req.Type),
	}
	if err := validation.CreateProductValidation(serviceReq); err != nil {
		return nil, err
	}
	return serviceReq, nil
}

func ParseListProductRequest(ctx kit.Context) (model.ListProductRequest, error) {
	q := ctx.Request().URL.Query()
	return model.ListProductRequest{
		Search: q.Get("search"),
		Type:   q.Get("type"),
		SortBy: q.Get("sort_by"),
		Order:  q.Get("order"),
	}, nil
}

func MapProductResponse(product *model.Product) model.CreateProductResponse {
	return model.CreateProductResponse{
		Product: product,
		Message: "Success",
	}
}

func MapListProductsResponse(products []model.Product) model.ListProductResponse {
	return model.ListProductResponse{
		Products: products,
		Message:  "Success",
	}
}
