package service

import (
	"context"

	pb "github.com/jlorgal/grpc-golab-19/book/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service type
type Service struct {
	log   *zap.Logger
	books []*pb.Book
}

// NewService creates a new book service
func NewService(log *zap.Logger) *Service {
	books := initBooks()
	log.Info("Initialized books", zap.Any("books", books))
	return &Service{
		log:   log,
		books: books,
	}
}

// GetBook retrieve a book given a bookID
func (s *Service) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	bookID := req.GetId()
	s.log.Info("Getting book", zap.Int64("bookID", bookID))
	book, err := s.findBook(bookID)
	if err != nil {
		s.log.Error("Error getting book", zap.Int64("bookID", bookID), zap.Error(err))
		return nil, err
	}
	s.log.Info("Got book", zap.Any("book", book))
	return &pb.GetBookResponse{Book: book}, nil
}

// GetAuthorBooks retrieve a list of books belonging to an author identified with authorID
func (s *Service) GetAuthorBooks(ctx context.Context, req *pb.GetAuthorBooksRequest) (*pb.GetAuthorBooksResponse, error) {
	authorID := req.GetAuthorId()
	s.log.Info("Getting author's books", zap.Int64("authorID", authorID))
	books := s.findAuthorBooks(authorID)
	s.log.Info("Got author's books", zap.Any("books", books))
	return &pb.GetAuthorBooksResponse{
		Books: books,
	}, nil
}

func (s *Service) findAuthorBooks(authorID int64) []*pb.Book {
	books := make([]*pb.Book, 0)
	for _, book := range s.books {
		if book.GetAuthorId() == authorID {
			books = append(books, book)
		}
	}
	return books
}

func (s *Service) findBook(bookID int64) (*pb.Book, error) {
	for _, book := range s.books {
		if book.GetId() == bookID {
			return book, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "Not found book %d", bookID)
}

func initBooks() []*pb.Book {
	return []*pb.Book{
		&pb.Book{Id: 1, Name: "Don Quixote", AuthorId: 1},
		&pb.Book{Id: 2, Name: "Rinconete and Cortadillo", AuthorId: 1},
		&pb.Book{Id: 3, Name: "Exemplary novels", AuthorId: 1},
		&pb.Book{Id: 4, Name: "Hamlet", AuthorId: 2},
		&pb.Book{Id: 5, Name: "Macbeth", AuthorId: 2},
		&pb.Book{Id: 6, Name: "Othello", AuthorId: 2},
		&pb.Book{Id: 7, Name: "Romeo and Juliet", AuthorId: 2},
	}
}
