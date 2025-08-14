package handler

import (
	"net/http"
	"strings"

	tpb "todo-gunk/gunk/v1/todo"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Todo struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	IsCompleted bool   `db:"is_completed"`
}

func (todo *Todo) Validate() error {
	return validation.ValidateStruct(todo,
		validation.Field(&todo.Title, validation.Required.Error("Title field can not be empty."), validation.Length(3, 50).Error("Title field should have atleast 3 characters and atmost 50 characters")),
		validation.Field(&todo.Description, validation.Required.Error("Description field can not be empty."), validation.Length(3, 50).Error("Description field should have atleast 3 characters and atmost 50 characters")),
	)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	ErrorValue := map[string]string{}
	todo := Todo{}
	h.loadCreateForm(w, todo, ErrorValue)
}

func (h *Handler) storeTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var todo Todo
	err = h.decoder.Decode(&todo, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = todo.Validate()
	if err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			ErrorValue := make(map[string]string)
			for key, value := range vErrors {
				ErrorValue[strings.Title(key)] = value.Error()
			}
			h.loadCreateForm(w, todo, ErrorValue)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.tc.Create(r.Context(), &tpb.CreateTodoRequest{
		Todo: &tpb.Todo{
			Title:       todo.Title,
			Description: todo.Description,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type FormData struct {
	Todo   Todo
	Errors map[string]string
}

func (h *Handler) loadCreateForm(w http.ResponseWriter, todo Todo, myErrors map[string]string) {
	form := FormData{
		Todo:   todo,
		Errors: myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "create-todo.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
