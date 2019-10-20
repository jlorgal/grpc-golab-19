package main

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	pb "github.com/jlorgal/grpc-golab-19/book/proto"
	"github.com/jlorgal/grpc-golab-19/book/service"
	"github.com/jlorgal/grpc-golab-19/metrics"
	"github.com/jlorgal/grpc-golab-19/tracer"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type configuration struct {
	Address             string `default:":4041"`
	TracerReporterURL   string `default:"http://localhost:9411/api/v2/spans"`
	TracerServiceName   string `default:"book"`
	TracerServiceTarget string `default:"localhost:4041"`
	MetricsTarget       string `default:":9999"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("book", &conf); err != nil {
		log.Fatal("Error processing book configuration", zap.Error(err))
	}

	svc := service.NewService(log)

	zipkinTracer, err := tracer.NewZipkinTracer(conf.TracerReporterURL, conf.TracerServiceName, conf.TracerServiceTarget)
	if err != nil {
		log.Fatal("Error configuring zipkin tracer", zap.Error(err))
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(zipkinTracer)),
		)),
	)

	pb.RegisterBookServiceServer(srv, svc)
	reflection.Register(srv)

	log.Info(fmt.Sprintf("Starting book service on: %s", conf.Address))
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		panic(err)
	}

	go metrics.ServePrometheusEndpoint(srv, conf.MetricsTarget)

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
