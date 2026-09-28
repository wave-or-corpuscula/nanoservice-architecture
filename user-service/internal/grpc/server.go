package grpcserver

import (
	"context"
	"errors"
	user "microservices/proto/user"
	"userservice/internal/database"
	"userservice/internal/handlers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	user.UnimplementedUserServiceServer
	db handlers.Database
}

func NewServer(db handlers.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) GetUser(
	ctx context.Context,
	req *user.GetUserRequest,
) (*user.GetUserResponse, error) {
	u, err := s.db.GetUser(uint(req.Id))
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &user.GetUserResponse{
		Id:    uint64(u.ID),
		Name:  u.Name,
		Email: u.Email,
	}, nil
}
