package todo

import (
	"context"

	tpb "todo-gunk/gunk/v1/todo"
	"todo-gunk/todo/storage"
)

type todoCoreStore interface {
	Create(context.Context, storage.Todo) (int64, error)
	Get(context.Context, int64) (*storage.Todo, error)
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
