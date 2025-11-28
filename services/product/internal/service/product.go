package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/yourusername/ecom/common/api" // Import generated code
	"github.com/yourusername/ecom/product/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductService implements the gRPC interface
type ProductService struct {
	Collection                            *mongo.Collection
	api.UnimplementedProductServiceServer // Required for forward compatibility
}

func (s *ProductService) CreateProduct(ctx context.Context, req *api.CreateProductRequest) (*api.CreateProductResponse, error) {
	// 1. Map Request to Model
	// Note: In real production, use a snowflake ID generator.
	// For now, we use a random Int64 to satisfy the proto contract.
	resID := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()

	product := models.Product{
		ID:          resID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	// 2. Insert into MongoDB
	_, err := s.Collection.InsertOne(ctx, product)
	if err != nil {
		return &api.CreateProductResponse{Error: "Failed to create product"}, status.Errorf(codes.Internal, "DB Error: %v", err)
	}

	return &api.CreateProductResponse{Id: resID}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *api.GetProductRequest) (*api.GetProductResponse, error) {
	var product models.Product

	// 1. Find by ID
	res := s.Collection.FindOne(ctx, bson.M{"_id": req.Id})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return &api.GetProductResponse{Error: "Product not found"}, status.Errorf(codes.NotFound, "Product not found")
		}
		return nil, status.Errorf(codes.Internal, "DB Error: %v", res.Err())
	}

	// 2. Decode into struct
	if err := res.Decode(&product); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to decode product")
	}

	// 3. Map to Proto Response
	return &api.GetProductResponse{
		Product: &api.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		},
	}, nil
}

// ListProducts Implementation
func (s *ProductService) ListProducts(ctx context.Context, req *api.ListProductsRequest) (*api.ListProductsResponse, error) {
	var products []*api.Product

	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch products")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p models.Product
		if err := cursor.Decode(&p); err != nil {
			continue
		}
		products = append(products, &api.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}

	return &api.ListProductsResponse{Products: products}, nil
}
