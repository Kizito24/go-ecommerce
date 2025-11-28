package main

import (
	"log"
	"net"

	"github.com/yourusername/ecom/common/api"
	"github.com/yourusername/ecom/product/internal/db"
	"github.com/yourusername/ecom/product/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// 1. Initialize DB
	mongoClient := db.Init()
	// Create/Get the collection "products" inside database "ecom_db"
	collection := mongoClient.Database("ecom_db").Collection("products")

	// 2. Setup gRPC Server
	// Listen on port 5002 (Matches your docker-compose)
	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// 3. Register the Product Service
	productService := &service.ProductService{
		Collection: collection,
	}
	api.RegisterProductServiceServer(grpcServer, productService)

	log.Println("Product Service listening on port 5002...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
