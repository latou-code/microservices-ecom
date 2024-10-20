package main

import (
	"log"
	"net"
	"order-service/client"
	"order-service/handler"
	pb "order-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Создание gRPC клиента для Inventory Service
	inventoryClient, err := client.NewInventoryClient("inventory-service:50056")
	if err != nil {
		log.Fatalf("Failed to connect to Inventory Service: %v", err)
	}

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, handler.NewOrderServiceServer(inventoryClient))

	log.Println("Order Service is running on port 50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
