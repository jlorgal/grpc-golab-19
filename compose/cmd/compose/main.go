package main

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	author_pb "github.com/jlorgal/grpc-golab-19/author/proto"
	book_pb "github.com/jlorgal/grpc-golab-19/book/proto"
	compose_pb "github.com/jlorgal/grpc-golab-19/compose/proto"
	"github.com/jlorgal/grpc-golab-19/compose/service"
	"github.com/jlorgal/grpc-golab-19/metrics"
	"github.com/jlorgal/grpc-golab-19/tracer"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type configuration struct {
	Address             string `default:":4042"`
	AuthorTarget        string `default:"127.0.0.1:4040"`
	BookTarget          string `default:"127.0.0.1:4041"`
	TracerReporterURL   string `default:"http://localhost:9411/api/v2/spans"`
	TracerServiceName   string `default:"compose"`
	TracerServiceTarget string `default:"localhost:4042"`
	MetricsTarget       string `default:":9999"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("compose", &conf); err != nil {
		log.Fatal("Error processing compose configuration", zap.Error(err))
	}

	zipkinTracer, err := tracer.NewZipkinTracer(conf.TracerReporterURL, conf.TracerServiceName, conf.TracerServiceTarget)
	if err != nil {
		log.Fatal("Error configuring zipkin tracer", zap.Error(err))
	}

	log.Info(fmt.Sprintf("Creating author client to: %s", conf.AuthorTarget))
	conn, err := grpc.Dial(conf.AuthorTarget,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(zipkinTracer)),
		),
	)
	if err != nil {
		log.Fatal("No connection to author service", zap.Error(err))
	}
	authorClient := author_pb.NewAuthorServiceClient(conn)

	log.Info(fmt.Sprintf("Creating book client to: %s", conf.BookTarget))
	conn2, err := grpc.Dial(conf.BookTarget,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(zipkinTracer)),
		),
	)
	if err != nil {
		log.Fatal("No connection to book service", zap.Error(err))
	}
	bookClient := book_pb.NewBookServiceClient(conn2)

	svc := service.NewService(log, authorClient, bookClient)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(zipkinTracer)),
		)),
	)

	compose_pb.RegisterComposeServiceServer(srv, svc)
	reflection.Register(srv)

	log.Info(fmt.Sprintf("Starting compose service on: %s", conf.Address))
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		panic(err)
	}

	go metrics.ServePrometheusEndpoint(srv, conf.MetricsTarget)

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
