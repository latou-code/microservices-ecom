package handler

import (
	"context"
	"log"
	"order-service/client"
	pb "order-service/proto"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	inventoryClient *client.InventoryClient
}

func NewOrderServiceServer(inventoryClient *client.InventoryClient) *OrderServiceServer {
	return &OrderServiceServer{inventoryClient: inventoryClient}
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Проверяем наличие всех товаров на складе
	for _, item := range req.Items {
		inStock, availableQuantity, err := s.inventoryClient.CheckStock(item.ProductId)
		if err != nil {
			log.Printf("Error checking stock for product %d: %v", item.ProductId, err)
			return &pb.CreateOrderResponse{Success: false, Message: "Error checking stock"}, err
		}

		if !inStock || availableQuantity < item.Quantity {
			log.Printf("Insufficient stock for product %d", item.ProductId)
			return &pb.CreateOrderResponse{Success: false, Message: "Insufficient stock"}, nil
		}
	}

	// Если все товары есть на складе, уменьшаем количество
	for _, item := range req.Items {
		success, err := s.inventoryClient.DecreaseStock(item.ProductId, item.Quantity)
		if err != nil || !success {
			log.Printf("Error decreasing stock for product %d: %v", item.ProductId, err)
			return &pb.CreateOrderResponse{Success: false, Message: "Error decreasing stock"}, err
		}
	}

	// Логика создания заказа (упрощенная)
	orderID := 1 // Допустим, что ID заказа 1

	return &pb.CreateOrderResponse{
		Success: true,
		Message: "Order created successfully",
		OrderId: int32(orderID),
	}, nil
}
