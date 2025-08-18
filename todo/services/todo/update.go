package todo

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tpb "todo-gunk/gunk/v1/todo"
	"todo-gunk/todo/storage"
)

func (s *Svc) Update(ctx context.Context, req *tpb.UpdateTodoRequest) (*tpb.UpdateTodoResponse, error) {
	_, err := s.core.Update(ctx, storage.Todo{
		ID:          req.Todo.ID,
		Title:       req.Todo.Title,
		Description: req.Todo.Description,
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &tpb.UpdateTodoResponse{}, nil
}
