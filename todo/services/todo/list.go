package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) List(ctx context.Context, req *tpb.ListTodoRequest) (*tpb.ListTodoResponse, error) {
	todos, err := s.core.List(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, "no todos found")
	}

	list := make([]*tpb.Todo, len(todos))
	for i, u := range todos {
		list[i] = &tpb.Todo{
			ID:          u.ID,
			Title:       u.Title,
			Description: u.Description,
			IsCompleted: u.IsCompleted,
		}
	}

	return &tpb.ListTodoResponse{
		Todo: list,
	}, nil
}
