package service

import (
	"context"

	author_pb "github.com/jlorgal/grpc-golab-19/author/proto"
	book_pb "github.com/jlorgal/grpc-golab-19/book/proto"
	compose_pb "github.com/jlorgal/grpc-golab-19/compose/proto"
	"go.uber.org/zap"
)

// Service type
type Service struct {
	log          *zap.Logger
	authorClient author_pb.AuthorServiceClient
	bookClient   book_pb.BookServiceClient
}

// NewService creates a new compose service
func NewService(log *zap.Logger, authorClient author_pb.AuthorServiceClient, bookClient book_pb.BookServiceClient) *Service {
	return &Service{
		log:          log,
		authorClient: authorClient,
		bookClient:   bookClient,
	}
}

// GetAuthor retrieves the author information from the author service and the
// author's books from the book service.
func (s *Service) GetAuthor(ctx context.Context, req *compose_pb.GetAuthorRequest) (*compose_pb.GetAuthorResponse, error) {
	authorID := req.GetId()
	s.log.Info("Getting author", zap.Int64("authorID", authorID))
	authorInfo, err := s.authorClient.GetAuthor(ctx, &author_pb.GetAuthorRequest{Id: authorID})
	if err != nil {
		s.log.Error("Error getting author from author service", zap.Error(err))
		return nil, err
	}
	s.log.Info("Author got from author service", zap.Any("author", authorInfo))

	authorBooks, err := s.bookClient.GetAuthorBooks(ctx, &book_pb.GetAuthorBooksRequest{AuthorId: authorID})
	if err != nil {
		s.log.Error("Error getting author's books from book service", zap.Error(err))
		return nil, err
	}
	s.log.Info("Author's books got from book service", zap.Any("books", authorBooks))

	return &compose_pb.GetAuthorResponse{
		Id:    authorID,
		Name:  authorInfo.GetAuthor().GetName(),
		Books: mapBooks(authorBooks.GetBooks()),
	}, nil
}

func mapBook(book *book_pb.Book) *compose_pb.Book {
	return &compose_pb.Book{
		Id:   book.GetId(),
		Name: book.GetName(),
	}
}

func mapBooks(books []*book_pb.Book) []*compose_pb.Book {
	result := make([]*compose_pb.Book, len(books))
	for i, book := range books {
		result[i] = mapBook(book)
	}
	return result
}
