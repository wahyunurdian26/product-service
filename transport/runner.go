package transport

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/wahyunurdian26/product-service/contract/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func (g *GrpcServer) Run() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "6668" // product service grpc port
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // product service http port
	}

	// Create gRPC Server
	g.handler = grpc.NewServer()
	pb.RegisterProductServiceServer(g.handler, g)
	grpc_health_v1.RegisterHealthServer(g.handler, g.healthServer)
	reflection.Register(g.handler)

	// Set health status
	for name := range g.handler.GetServiceInfo() {
		g.healthServer.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	}

	// HTTP Gateway Setup
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Gateway dials the local gRPC server
	err := pb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}

	go g.waitForShutdown()

	log.Printf("Product Service started successfully - GRPC on :%s, REST Gateway on :%s\n", grpcPort, httpPort)

	// Start HTTP Server
	go func() {
		if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Start gRPC Server
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}
	if err := g.handler.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func (g *GrpcServer) waitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-c
	log.Println("Shutdown signal received")
	g.Stop()
}

func (g *GrpcServer) Stop() {
	log.Println("Stopping Server - initiating graceful shutdown")
	if g.handler != nil {
		g.handler.GracefulStop()
	}
	g.close()
	log.Println("Server stopped successfully")
}
