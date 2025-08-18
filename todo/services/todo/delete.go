package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"
)

func (s *Svc) Delete(ctx context.Context, req *tpb.DeleteTodoRequest) (*tpb.DeleteTodoResponse, error) {
	if err := s.core.Delete(ctx, req.GetID()); err != nil {
		return nil, err
	}

	return &tpb.DeleteTodoResponse{}, nil
}
