package http

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/wahyunurdian26/product-service/endpoint"
	"github.com/wahyunurdian26/product-service/model"
)

func NewHTTPHandler(endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/product").Handler(httptransport.NewServer(
		endpoints.CreateProduct,
		decodeCreateProductRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/product").Handler(httptransport.NewServer(
		endpoints.ListProducts,
		decodeListProductsRequest,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func decodeCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func decodeListProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	search := r.URL.Query().Get("search")
	productType := r.URL.Query().Get("type")
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")

	return &model.ListProductRequest{
		Search: search,
		Type:   productType,
		SortBy: sortBy,
		Order:  order,
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
