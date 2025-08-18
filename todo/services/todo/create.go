package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"
	"todo-gunk/todo/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Create(ctx context.Context, req *tpb.CreateTodoRequest) (*tpb.CreateTodoResponse, error) {
	//Need to validate request here

	todo := storage.Todo{
		ID:          req.GetTodo().ID,
		Title:       req.GetTodo().Title,
		Description: req.GetTodo().Description,
	}

	id, err := s.core.Create(context.Background(), todo)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create todo")
	}

	return &tpb.CreateTodoResponse{
		ID: id,
	}, nil
}
