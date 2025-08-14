package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Get(ctx context.Context, req *tpb.GetTodoRequest) (*tpb.GetTodoResponse, error) {

	id := req.GetID()
	todo, err := s.core.Get(context.Background(), id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get todo")
	}

	return &tpb.GetTodoResponse{
		Todo: &tpb.Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		},
	}, nil
}
