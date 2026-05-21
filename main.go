package main

import (
	_ "github.com/joho/godotenv/autoload"
	httptransport "github.com/wahyunurdian26/product-service/transport/http"
)

func main() {
	srv := httptransport.NewHTTPServer()
	srv.Run()
}
