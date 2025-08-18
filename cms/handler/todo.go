package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	tpb "todo-gunk/gunk/v1/todo"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
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

func (h *Handler) editTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.tc.Get(r.Context(), &tpb.GetTodoRequest{
		ID: int64(id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.loadEditForm(w, Todo{
		ID:          res.Todo.ID,
		Title:       res.Todo.Title,
		Description: res.Todo.Description,
	}, map[string]string{})
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var todo Todo
	if err := h.decoder.Decode(&todo, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := todo.Validate(); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for k, v := range e {
					vErrs[k] = v.Error()
				}
			}
		}
		h.loadEditForm(w, todo, vErrs)
		return
	}

	if _, err := h.tc.Update(r.Context(), &tpb.UpdateTodoRequest{
		Todo: &tpb.Todo{
			ID:          int64(id),
			Title:       todo.Title,
			Description: todo.Description,
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) listTodo(w http.ResponseWriter, r *http.Request) {
	template := h.templates.Lookup("list-todo.html")
	if template == nil {
		errMsg := "unable to load template"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	todos, err := h.tc.List(r.Context(), &tpb.ListTodoRequest{})
	if err != nil {
		log.Println("unable to get list: ", err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	ltList := make([]Todo, 0, len(todos.GetTodo()))
	for _, item := range todos.GetTodo() {
		listData := Todo{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			IsCompleted: item.IsCompleted,
		}
		ltList = append(ltList, listData)
	}

	if err := template.Execute(w, ltList); err != nil {
		log.Printf("error with template execution: %+v", err)
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.tc.Delete(r.Context(), &tpb.DeleteTodoRequest{
		ID: int64(id),
	}); err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) completeTodo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := h.tc.Complete(r.Context(), &tpb.CompleteTodoRequest{
		ID: int64(id),
	}); err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (h *Handler) loadEditForm(w http.ResponseWriter, todo Todo, myErrors map[string]string) {
	form := FormData{
		Todo:   todo,
		Errors: myErrors,
	}

	err := h.templates.ExecuteTemplate(w, "edit-todo.html", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
