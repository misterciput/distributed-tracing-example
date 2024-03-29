// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: shipping/shipping.proto

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	shipping/shipping.proto

It has these top-level messages:
	Request
	Response
	Shipping
*/
package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ShippingService service

type ShippingService interface {
	GetShippingProvider(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type shippingService struct {
	c    client.Client
	name string
}

func NewShippingService(name string, c client.Client) ShippingService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "main"
	}
	return &shippingService{
		c:    c,
		name: name,
	}
}

func (c *shippingService) GetShippingProvider(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "ShippingService.GetShippingProvider", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ShippingService service

type ShippingServiceHandler interface {
	GetShippingProvider(context.Context, *Request, *Response) error
}

func RegisterShippingServiceHandler(s server.Server, hdlr ShippingServiceHandler, opts ...server.HandlerOption) error {
	type shippingService interface {
		GetShippingProvider(ctx context.Context, in *Request, out *Response) error
	}
	type ShippingService struct {
		shippingService
	}
	h := &shippingServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ShippingService{h}, opts...))
}

type shippingServiceHandler struct {
	ShippingServiceHandler
}

func (h *shippingServiceHandler) GetShippingProvider(ctx context.Context, in *Request, out *Response) error {
	return h.ShippingServiceHandler.GetShippingProvider(ctx, in, out)
}
