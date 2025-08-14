package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"todo-gunk/cms/handler"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	tpb "todo-gunk/gunk/v1/todo"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.GetString("todo.host"), config.GetString("todo.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Connection failed", err)
	}

	tc := tpb.NewTodoServiceClient(conn)
	r := handler.GetHandler(decoder, store, tc)

	host, port := config.GetString("server.host"), config.GetString("server.port")
	log.Println("Server  starting...")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
