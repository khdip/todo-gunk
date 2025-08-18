package todo

import (
	"context"

	"todo-gunk/todo/storage"
)

type todoStore interface {
	Create(ctx context.Context, t storage.Todo) (int64, error)
	Get(ctx context.Context, id int64) (*storage.Todo, error)
	Update(ctx context.Context, t storage.Todo) (*storage.Todo, error)
	List(ctx context.Context) ([]storage.Todo, error)
	Delete(ctx context.Context, id int64) error
	Complete(ctx context.Context, id int64) error
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

func (cs CoreSvc) Update(ctx context.Context, t storage.Todo) (*storage.Todo, error) {
	return cs.store.Update(ctx, t)
}

func (cs CoreSvc) List(ctx context.Context) ([]storage.Todo, error) {
	return cs.store.List(ctx)
}

func (cs CoreSvc) Delete(ctx context.Context, id int64) error {
	return cs.store.Delete(ctx, id)
}

func (cs CoreSvc) Complete(ctx context.Context, id int64) error {
	return cs.store.Complete(ctx, id)
}
