package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"

	tpb "todo-gunk/gunk/v1/todo"
)

// const sessionName = "cms-session"

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
	// r.HandleFunc("/", hand.home)

	s := r.NewRoute().Subrouter()
	s.HandleFunc("/create", hand.createTodo)
	s.HandleFunc("/store", hand.storeTodo)
	// s.Use(hand.authMiddleware)
	// r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
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
		"cms/assets/templates/404.html",
	))
}
