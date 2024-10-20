package client

import (
	"context"
	"log"
	"time"

	pb "order-service/proto/inventory" // Подключаем сгенерированный код Inventory
	"google.golang.org/grpc"
)

type InventoryClient struct {
	client pb.InventoryServiceClient
}

func NewInventoryClient(address string) (*InventoryClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*5))
	if err != nil {
		return nil, err
	}

	c := pb.NewInventoryServiceClient(conn)
	return &InventoryClient{client: c}, nil
}

func (c *InventoryClient) CheckStock(productID int32) (bool, int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.client.CheckStock(ctx, &pb.CheckStockRequest{ProductId: productID})
	if err != nil {
		return false, 0, err
	}

	return resp.InStock, resp.AvailableQuantity, nil
}

func (c *InventoryClient) DecreaseStock(productID int32, quantity int32) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.client.DecreaseStock(ctx, &pb.DecreaseStockRequest{ProductId: productID, Quantity: quantity})
	if err != nil {
		return false, err
	}

	return resp.Success, nil
}
