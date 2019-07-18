package main

import (
	"context"
	"log"
	"net"

	tracer "github.com/misterciput/meetup/tracer"
	opentracing "github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
)

type (
	serverGrpc struct{}
)

const (
	port = ":50052"
)

func init() {
	tracer.InitJaeger("product-service")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterProductServiceServer(s, &serverGrpc{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *serverGrpc) GetListProduct(ctx context.Context, in *ProductRequest) (*ProductResponse, error) {
	span, ctx := tracer.StartSpanWithExtract(ctx, "GetShippingProvider")
	defer span.Finish()

	response := getProductDetail(ctx)

	for key, _ := range response.Products {
		response.Products[key].Price = getProductPrice(ctx, response.Products[key].ProductID)
	}

	return response, nil
}

func getProductDetail(ctx context.Context) *ProductResponse {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductDetail")
	defer span.Finish()

	response := &ProductResponse{
		Products: []*Product{
			{
				ProductID:   10,
				Name:        "Samsul Galaxy S10",
				Description: "Barang Palsu, silahkan diorder",
			},
			{
				ProductID:   20,
				Name:        "aiPhone X",
				Description: "Barang KW, silahkan diorder",
			},
		},
	}

	return response
}

func getProductPrice(ctx context.Context, productID int32) int32 {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductPrice")
	defer span.Finish()

	switch productID {
	case 10:
		return 100000
	case 20:
		return 200000
	}

	return 0
}
