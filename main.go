package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/wahyunurdian26/product-service/transport"
)

func main() {
	srv := transport.NewGRPCServer()
	srv.Run()
}
