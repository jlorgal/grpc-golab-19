package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	compose_pb "github.com/jlorgal/grpc-golab-19/compose/proto"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type configuration struct {
	Address       string `default:":8080"`
	ComposeTarget string `default:"127.0.0.1:4043"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("api", &conf); err != nil {
		log.Fatal("Error processing api configuration", zap.Error(err))
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := compose_pb.RegisterComposeServiceHandlerFromEndpoint(ctx, mux, conf.ComposeTarget, opts); err != nil {
		log.Fatal("Error registering api endpoint handler", zap.Error(err))
	}

	log.Info(fmt.Sprintf("Starting api service on: %s", conf.Address))
	if err := http.ListenAndServe(conf.Address, mux); err != nil {
		log.Fatal("Error serving api service", zap.Error(err))
	}
}
