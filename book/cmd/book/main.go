package main

import (
	"fmt"
	"net"

	pb "github.com/jlorgal/grpc-golab-19/book/proto"
	"github.com/jlorgal/grpc-golab-19/book/service"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type configuration struct {
	Address string `default:":4041"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("book", &conf); err != nil {
		log.Fatal("Error processing book configuration", zap.Error(err))
	}

	svc := service.NewService(log)
	srv := grpc.NewServer()
	pb.RegisterBookServiceServer(srv, svc)
	reflection.Register(srv)

	log.Info(fmt.Sprintf("Starting book service on: %s", conf.Address))
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		panic(err)
	}

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
