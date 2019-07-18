package tracer

import (
	"context"
	"fmt"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	metadataGrpc "google.golang.org/grpc/metadata"
)

func InitJaeger(appName string) {

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return
	}

	jLogger := jaegerlog.NullLogger // use StdLogger if need to debug
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	_, err = cfg.InitGlobalTracer(
		appName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
}

// Create new context based on span
// Use it when you want to send context between service
func Inject(ctx context.Context, span opentracing.Span) context.Context {

	spanContext := span.Context()
	jaegerSpanContext, ok := spanContext.(jaeger.SpanContext)
	if !ok {
		return ctx
	}

	spanContextStr := jaegerSpanContext.String()

	md, _ := metadataGrpc.FromIncomingContext(ctx)

	md2 := make(map[string][]string)
	md2["uber-trace-id"] = make([]string, 0)
	md2["uber-trace-id"] = append(md2["uber-trace-id"], spanContextStr)

	for key, val := range md {
		fmt.Println("KEY: ", key, "| VALUE: ", val)
		if key != "uber-trace-id" {
			md2[key] = val
		}
	}

	// create new context based on span and updated SPAN CONTEXT
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadataGrpc.NewOutgoingContext(ctx, md2)

	return ctx
}

// Start new span based on context
// Use it when you want to start span by receiving context from another service
func StartSpanWithExtract(ctx context.Context, name string) (opentracing.Span, context.Context) {

	md, _ := metadataGrpc.FromIncomingContext(ctx)

	mdCarrier := make(map[string]string)
	for key, val := range md {
		fmt.Println("KEY: ", key, "| VALUE: ", val)
		mdCarrier[key] = val[0]
	}

	// extract SPAN CONTEXT from given context
	wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(mdCarrier))
	if err != nil {
		return opentracing.StartSpan(name), ctx
	}

	// start new child span based on previous context
	span := opentracing.StartSpan(name, opentracing.ChildOf(wireContext))
	ctx = Inject(ctx, span)

	return span, ctx
}
