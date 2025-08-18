package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	tpb "todo-gunk/gunk/v1/todo"
)

type Handler struct {
	templates *template.Template
	decoder   *schema.Decoder
	session   *sessions.CookieStore
	tc        tpb.TodoServiceClient
}

func GetHandler(decoder *schema.Decoder, session *sessions.CookieStore, tc tpb.TodoServiceClient) *mux.Router {
	hand := &Handler{
		decoder: decoder,
		session: session,
		tc:      tc,
	}
	hand.GetTemplate()

	r := mux.NewRouter()
	r.HandleFunc("/", hand.listTodo)
	r.HandleFunc("/create", hand.createTodo)
	r.HandleFunc("/store", hand.storeTodo)
	r.HandleFunc("/edit/{id}", hand.editTodo)
	r.HandleFunc("/update/{id}", hand.updateTodo)
	r.HandleFunc("/delete/{id}", hand.deleteTodo)
	r.HandleFunc("/complete/{id}", hand.completeTodo)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := hand.templates.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) GetTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"cms/assets/templates/create-todo.html",
		"cms/assets/templates/edit-todo.html",
		"cms/assets/templates/404.html",
		"cms/assets/templates/list-todo.html",
	))
}
