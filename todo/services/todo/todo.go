package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"
	"todo-gunk/todo/storage"
)

type todoCoreStore interface {
	Create(context.Context, storage.Todo) (int64, error)
	Get(context.Context, int64) (*storage.Todo, error)
	Update(context.Context, storage.Todo) (*storage.Todo, error)
	List(context.Context) ([]storage.Todo, error)
	Delete(context.Context, int64) error
	Complete(context.Context, int64) error
}

type Svc struct {
	tpb.UnimplementedTodoServiceServer
	core todoCoreStore
}

func NewTodoServer(c todoCoreStore) *Svc {
	return &Svc{
		core: c,
	}
}
