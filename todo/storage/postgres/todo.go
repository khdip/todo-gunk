package postgres

import (
	"context"
	"fmt"
	"todo-gunk/todo/storage"
)

const insertTodo = `
	INSERT INTO todos(
		title,
		description
	) VALUES(
		:title,
		:description
	)
	RETURNING id;
`

func (s *Storage) Create(ctx context.Context, t storage.Todo) (int64, error) {
	stmt, err := s.db.PrepareNamed(insertTodo)
	if err != nil {
		return 0, err
	}

	var id int64
	if err := stmt.Get(&id, t); err != nil {
		return 0, err
	}
	return id, nil
}

const getTodo = `
SELECT 
	id,
	title, 
	description,
	is_completed
FROM todos
WHERE id = $1;
`

func (s *Storage) Get(ctx context.Context, id int64) (*storage.Todo, error) {
	var todo storage.Todo
	if err := s.db.Get(&todo, getTodo, id); err != nil {
		return nil, fmt.Errorf("executing leave_type details: %w", err)
	}
	return &todo, nil
}
