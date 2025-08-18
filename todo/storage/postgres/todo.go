package postgres

import (
	"context"
	"database/sql"
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

const updateTodo = `
UPDATE
	todos
SET
	title = :title,
    description = :description
WHERE 
	id = :id
RETURNING *;
`

func (s *Storage) Update(ctx context.Context, todo storage.Todo) (*storage.Todo, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateTodo)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.Get(&todo, todo); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found: %w", err)
		}
		return nil, fmt.Errorf("executing todo update: %w", err)
	}

	return &todo, nil
}

const listTodo = `
SELECT * FROM todos;
`

func (s *Storage) List(ctx context.Context) ([]storage.Todo, error) {
	var todos []storage.Todo
	if err := s.db.Select(&todos, listTodo); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("not found: %w", err)
		}
		return nil, err
	}

	return todos, nil
}

const deleteTodo = `
DELETE FROM todos WHERE id = $1;
`

func (s *Storage) Delete(ctx context.Context, id int64) error {
	_, err := s.db.Exec(deleteTodo, id)
	if err != nil {
		panic(err)
	}
	return nil
}

const completeTodo = `
UPDATE todos SET is_completed = true WHERE id = $1;
`

func (s *Storage) Complete(ctx context.Context, id int64) error {
	_, err := s.db.Exec(completeTodo, id)
	if err != nil {
		panic(err)
	}
	return nil
}
