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
	port = ":50051"
)

func init() {
	tracer.InitJaeger("shipping-service")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterShippingServiceServer(s, &serverGrpc{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *serverGrpc) GetShippingProvider(ctx context.Context, in *Request) (*Response, error) {
	span, ctx := tracer.StartSpanWithExtract(ctx, "GetShippingProvider")
	defer span.Finish()

	response := getShippingDetail(ctx)

	for key, _ := range response.Shippings {
		response.Shippings[key].PricePerUnit = getShippingPrice(ctx, response.Shippings[key].ShippingID)
	}

	return response, nil
}

func getShippingDetail(ctx context.Context) *Response {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getShippingDetail")
	defer span.Finish()

	response := &Response{
		Shippings: []*Shipping{
			{
				ShippingID: 10,
				Name:       "JNX",
				Unit:       "KG",
			},
			{
				ShippingID: 20,
				Name:       "GREP SEND",
				Unit:       "KM",
			},
		},
	}

	return response
}

func getShippingPrice(ctx context.Context, shippingID int32) int32 {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getShippingPrice")
	defer span.Finish()

	switch shippingID {
	case 10:
		return 8000
	case 20:
		return 10000
	}

	return 0
}
