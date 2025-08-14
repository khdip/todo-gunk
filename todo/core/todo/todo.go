package todo

import (
	"context"

	"todo-gunk/todo/storage"
)

type todoStore interface {
	Create(ctx context.Context, t storage.Todo) (int64, error)
	Get(ctx context.Context, id int64) (*storage.Todo, error)
}

type CoreSvc struct {
	store todoStore
}

func NewCoreSvc(s todoStore) *CoreSvc {
	return &CoreSvc{
		store: s,
	}
}

func (cs CoreSvc) Create(ctx context.Context, t storage.Todo) (int64, error) {
	return cs.store.Create(ctx, t)
}

func (cs CoreSvc) Get(ctx context.Context, id int64) (*storage.Todo, error) {
	return cs.store.Get(ctx, id)
}
