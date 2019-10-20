package service

import (
	"context"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/jlorgal/grpc-golab-19/author/proto"
)

// Service type
type Service struct {
	log     *zap.Logger
	authors []*pb.Author
}

// NewService creates a new author service
func NewService(log *zap.Logger) *Service {
	authors := initAuthors()
	log.Info("Initialized authors", zap.Any("authors", authors))
	return &Service{
		log:     log,
		authors: authors,
	}
}

// GetAuthor retrieves the author information given an authorID
func (s *Service) GetAuthor(ctx context.Context, req *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
	authorID := req.GetId()
	s.log.Info("Getting author", zap.Int64("authorID", authorID))
	author, err := s.findAuthor(authorID)
	if err != nil {
		s.log.Error("Error getting author", zap.Int64("authorID", authorID), zap.Error(err))
		return nil, err
	}
	s.log.Info("Got author", zap.Any("author", author))
	return &pb.GetAuthorResponse{Author: author}, nil
}

func (s *Service) findAuthor(authorID int64) (*pb.Author, error) {
	for _, author := range s.authors {
		log.Printf("checking author: %+v", author)
		if author.GetId() == authorID {
			return author, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Not found author %d", authorID)
}

func initAuthors() []*pb.Author {
	return []*pb.Author{
		&pb.Author{Id: 1, Name: "Miguel de Cervantes"},
		&pb.Author{Id: 2, Name: "William Shakespeare"},
	}
}
