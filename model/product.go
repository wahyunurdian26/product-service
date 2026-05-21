package model

import "time"

// ProductType defines the enum for product types
type ProductType string

const (
	TypeSayuran ProductType = "Sayuran"
	TypeProtein ProductType = "Protein"
	TypeBuah    ProductType = "Buah"
	TypeSnack   ProductType = "Snack"
)

// Product represents the product entity in the database
type Product struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Price     float64     `json:"price"`
	Type      ProductType `json:"type"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// CreateProductRequest is the payload for adding a new product
type CreateProductRequest struct {
	Name  string      `json:"name"`
	Price float64     `json:"price"`
	Type  ProductType `json:"type"`
}

// CreateProductResponse is the response payload
type CreateProductResponse struct {
	Product *Product `json:"product,omitempty"`
	Message string   `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// ListProductRequest is the parameter for listing products
type ListProductRequest struct {
	Search string `json:"search"`
	Type   string `json:"type"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

// ListProductResponse is the response payload
type ListProductResponse struct {
	Products []Product `json:"products"`
	Message  string    `json:"message,omitempty"`
	Error    string    `json:"error,omitempty"`
}
