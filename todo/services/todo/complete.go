package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"
)

func (s *Svc) Complete(ctx context.Context, req *tpb.CompleteTodoRequest) (*tpb.CompleteTodoResponse, error) {
	if err := s.core.Complete(ctx, req.GetID()); err != nil {
		return nil, err
	}

	return &tpb.CompleteTodoResponse{}, nil
}
