package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/mozartted/simple_todo_server/internal/config"
	handlers "github.com/mozartted/simple_todo_server/internal/handlers"
)

type App struct {
	httpServer *http.Server
	db         *sql.DB
	handlers.RouteHandler
}

func newApp(db *sql.DB) *App {
	return &App{
		db: db,
		RouteHandler: handlers.RouteHandler{
			Db: db,
		},
	}
}

func (a *App) RunServer() {
	appServer := mux.NewRouter()
	if err := config.InitialDBMigration(a.db); err != nil {
		log.Fatal(err)
	}

	appServer.HandleFunc("/todo", a.TodoListGet).Methods("GET")
	appServer.HandleFunc("/todo", a.TodoListAdd).Methods("POST")
	appServer.HandleFunc("/todo/{todoId}", a.TodoListDelete).Methods("DELETE")
	appServer.HandleFunc("/todo/{todoId}/update", a.TodoListUpdate).Methods("POST")
	// appServer.HandleFunc("/todo", a.TodoListGet)
	a.httpServer = &http.Server{
		Addr:    ":8082",
		Handler: appServer,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		fmt.Println("Server running live.")
		fmt.Print("Running function\n")
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Println(err.Error())
			cancel()
		}
		cancel()
	}(cancel)

	cancelSignal := make(chan os.Signal, 1)

	signal.Notify(cancelSignal, os.Interrupt)

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			log.Fatalln(err)
		}
	case <-cancelSignal:
		defer cancel()
		a.httpServer.Shutdown(ctx)
		log.Println("Graceful shutdown in process")

	}
}

func main() {
	dbConnection, err := config.DbConnect()
	if err != nil {
		log.Fatal(err)
	}

	app := newApp(dbConnection)
	app.RunServer()
}
