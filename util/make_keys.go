package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wahyunurdian26/product-service/model"
)

const PrefixKey = "product-service"

func MakeListProductsKey(req model.ListProductRequest) string {
	reqBytes, _ := json.Marshal(req)
	hash := sha256.Sum256(reqBytes)
	return fmt.Sprintf("%s:products_list:%s", PrefixKey, hex.EncodeToString(hash[:]))
}

func MakeProductsPattern() string {
	return fmt.Sprintf("%s:products_list:*", PrefixKey)
}
