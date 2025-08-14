package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	tpb "todo-gunk/gunk/v1/todo"
	tc "todo-gunk/todo/core/todo"
	"todo-gunk/todo/services/todo"
	"todo-gunk/todo/storage/postgres"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("todo/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	grpcServer := grpc.NewServer()
	store, err := newDBFromConfig(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	cs := tc.NewCoreSvc(store)
	s := todo.NewTodoServer(cs)
	tpb.RegisterTodoServiceServer(grpcServer, s)

	host, port := config.GetString("server.host"), config.GetString("server.port")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	log.Printf("Server is starting on: %s:%s", host, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslmode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransactionTimeout")
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
