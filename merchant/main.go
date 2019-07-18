package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	tracer "github.com/misterciput/meetup/tracer"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type (
	MerchantDetail struct {
		Merchant  Merchant       `json:"merchant"`
		Shippings []ShippingData `json:"shippings"`
		Products  []ProductData  `json:"products"`
	}

	Merchant struct {
		MerchantID int    `json:"merchant_id"`
		Name       string `json:"name"`
		Address    string `json:"address"`
		Phone      string `json:"phone"`
	}

	ShippingData struct {
		ShippingID   int    `json:"shipping_id"`
		Name         string `json:"name"`
		Unit         string `json:"unit"`
		PricePerUnit int    `json:"price_per_unit"`
	}

	ProductData struct {
		ProductID   int    `json:"product_id"`
		Name        string `json:"name"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}
)

const (
	addressShipping = "localhost:50051"
	addressProduct  = "localhost:50052"
)

func init() {
	tracer.InitJaeger("merchant-service")
}

func main() {

	http.HandleFunc("/v1/detail", handlerMerchantDetail)

	http.ListenAndServe(":80", nil)
}

func handlerMerchantDetail(res http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	span, ctx := opentracing.StartSpanFromContext(ctx, "handlerMerchantDetail")
	defer span.Finish()

	merchant := getMerchant(ctx)
	shippings := getShippings(ctx)
	products := getProducts(ctx)

	merchantDetail := MerchantDetail{
		Merchant:  merchant,
		Shippings: shippings,
		Products:  products,
	}

	responseJSON, _ := json.Marshal(merchantDetail)

	res.Write(responseJSON)

}

func getMerchant(ctx context.Context) Merchant {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getMerchant")
	defer span.Finish()

	var merchant Merchant
	merchant.MerchantID = 100
	merchant.Name = "Mawar"
	merchant.Address = "Jl. Inajadulu No.7, Jakarta Selatan"
	merchant.Phone = "080989999"

	return merchant
}

func getShippings(ctx context.Context) []ShippingData {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getShippings")
	defer span.Finish()

	conn, err := grpc.Dial(addressShipping, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	shippingService := NewShippingServiceClient(conn)

	ctx = tracer.Inject(ctx, span)
	response, _ := shippingService.GetShippingProvider(ctx, &Request{})

	shippings := make([]ShippingData, 0)

	for _, shp := range response.Shippings {
		shippings = append(shippings, ShippingData{
			ShippingID:   int(shp.ShippingID),
			Name:         shp.Name,
			Unit:         shp.Unit,
			PricePerUnit: int(shp.PricePerUnit),
		})
	}

	return shippings
}

func getProducts(ctx context.Context) []ProductData {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProducts")
	defer span.Finish()

	conn, err := grpc.Dial(addressProduct, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	productService := NewProductServiceClient(conn)

	ctx = tracer.Inject(ctx, span)
	response, _ := productService.GetListProduct(ctx, &ProductRequest{})

	products := make([]ProductData, 0)

	for _, pd := range response.Products {
		products = append(products, ProductData{
			ProductID:   int(pd.ProductID),
			Name:        pd.Name,
			Price:       int(pd.Price),
			Description: pd.Description,
		})
	}

	return products
}
