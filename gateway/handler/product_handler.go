package handler


import (
	"github.com/wahyunurdian26/product-service/gateway/handler/dto"
	"github.com/wahyunurdian26/product-service/gateway/kit"
	"github.com/wahyunurdian26/product-service/service"
	uLog "github.com/wahyunurdian26/util/logger"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) InitializeRoute(r *kit.Router) {
	r.Post("/product", h.createProduct)
	r.Get("/product", h.listProducts)
}

func (h *ProductHandler) createProduct(ctx kit.Context) (interface{}, error) {
	req, err := dto.ParseCreateProductRequest(ctx)
	if err != nil {
		uLog.LogError(ctx, "dto.ParseCreateProductRequest", err)
		return nil, err
	}

	uLog.LogRequest(ctx, "createProduct", req)

	serviceReq, err := req.ToServiceRequest()
	if err != nil {
		uLog.LogError(ctx, "req.ToServiceRequest", err)
		return nil, err
	}

	product, err := h.svc.CreateProduct(ctx.Context, serviceReq)
	if err != nil {
		uLog.LogError(ctx, "h.svc.CreateProduct", err)
		return nil, err
	}

	return dto.MapProductResponse(product), nil
}

func (h *ProductHandler) listProducts(ctx kit.Context) (interface{}, error) {
	req, err := dto.ParseListProductRequest(ctx)
	if err != nil {
		uLog.LogError(ctx, "dto.ParseListProductRequest", err)
		return nil, err
	}

	uLog.LogRequest(ctx, "listProducts", req)

	products, err := h.svc.ListProducts(ctx.Context, req)
	if err != nil {
		uLog.LogError(ctx, "h.svc.ListProducts", err)
		return nil, err
	}

	return dto.MapListProductsResponse(products), nil
}
