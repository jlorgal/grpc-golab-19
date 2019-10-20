package main

import (
	"fmt"
	"net"

	author_pb "github.com/jlorgal/grpc-golab-19/author/proto"
	book_pb "github.com/jlorgal/grpc-golab-19/book/proto"
	compose_pb "github.com/jlorgal/grpc-golab-19/compose/proto"
	"github.com/jlorgal/grpc-golab-19/compose/service"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type configuration struct {
	Address      string `default:":4042"`
	AuthorTarget string `default:"127.0.0.1:4040"`
	BookTarget   string `default:"127.0.0.1:4041"`
}

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	var conf configuration
	if err := envconfig.Process("compose", &conf); err != nil {
		log.Fatal("Error processing compose configuration", zap.Error(err))
	}

	log.Info(fmt.Sprintf("Creating author client to: %s", conf.AuthorTarget))
	conn, err := grpc.Dial(conf.AuthorTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatal("No connection to author service", zap.Error(err))
	}
	authorClient := author_pb.NewAuthorServiceClient(conn)

	log.Info(fmt.Sprintf("Creating book client to: %s", conf.BookTarget))
	conn2, err := grpc.Dial(conf.BookTarget, grpc.WithInsecure())
	if err != nil {
		log.Fatal("No connection to book service", zap.Error(err))
	}
	bookClient := book_pb.NewBookServiceClient(conn2)

	svc := service.NewService(log, authorClient, bookClient)
	srv := grpc.NewServer()
	compose_pb.RegisterComposeServiceServer(srv, svc)
	reflection.Register(srv)

	log.Info(fmt.Sprintf("Starting compose service on: %s", conf.Address))
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		panic(err)
	}

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
