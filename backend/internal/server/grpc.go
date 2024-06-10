package server

import (
	"context"

	v1 "backend/api/helloworld/v1"
	"backend/internal/conf"
	"backend/internal/service"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	logger log.Logger,
	tr *conf.Trace,
	greeter *service.GreeterService,
) *grpc.Server {
	// trace start
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			// serviceName,
			semconv.ServiceNameKey.String(tr.Jaeger.ServiceName),
			attribute.String("exporter", "otlp trace grpc"),
			attribute.String("environment", "dev"),
			attribute.Float64("float", 312.23),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// shutdownTracerProvider, err := initTracerProvider(ctx, res, tr.Jaeger.Http.Endpoint)
	_, err2 := initTracerProvider(ctx, res, tr.Jaeger.Grpc.Endpoint)
	if err2 != nil {
		log.Fatal(err)
	}
	// trace end
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(), // trace 链路追踪
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServiceServer(srv, greeter)
	return srv
}
