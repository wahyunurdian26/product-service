package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wahyunurdian26/product-service/model"
)

const PrefixKey = "product-service"

// MakeListProductsKey generates a unique cache key based on the list product request filters
func MakeListProductsKey(req model.ListProductRequest) string {
	reqBytes, _ := json.Marshal(req)
	hash := sha256.Sum256(reqBytes)
	return fmt.Sprintf("%s:products_list:%s", PrefixKey, hex.EncodeToString(hash[:]))
}

// MakeProductsPattern generates the pattern to invalidate all products list caches
func MakeProductsPattern() string {
	return fmt.Sprintf("%s:products_list:*", PrefixKey)
}
