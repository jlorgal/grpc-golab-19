package main

import (
	"fmt"
	"net"

	pb "github.com/jlorgal/grpc-golab-19/author/proto"
	"github.com/jlorgal/grpc-golab-19/author/service"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type configuration struct {
	Address string `default:":4040"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("author", &conf); err != nil {
		log.Fatal("Error processing author configuration", zap.Error(err))
	}

	svc := service.NewService(log)
	srv := grpc.NewServer()
	pb.RegisterAuthorServiceServer(srv, svc)
	reflection.Register(srv)

	log.Info(fmt.Sprintf("Starting author service on: %s", conf.Address))
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		panic(err)
	}

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
