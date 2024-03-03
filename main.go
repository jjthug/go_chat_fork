package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"server/api"
	db "server/db/sqlc"
	"server/util"
	"server/ws"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load config")
		log.Fatal("cannot load config")
	}

	fmt.Println("loaded config")

	// Postgres connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	fmt.Println("connected to pg db")

	store := db.NewStore(connPool)

	// websocket handler
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	server, err := api.NewServer(config, store, wsHandler)

	if err != nil {
		log.Fatal("error creating server", err)
	}
	err = server.RunHTTPServer(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
