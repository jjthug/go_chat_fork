package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"server/api"
	db "server/db/sqlc"
	"server/util"
	"server/ws_worker"
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

	// todo db migrations

	//todo redis server connection

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServer,
	}

	broadcastChannel := make(chan *ws_worker.Message, 5)
	taskDistributor := ws_worker.NewRedisTaskDistributor(redisOpt)
	taskProcessor := ws_worker.NewRedisTaskProcessor(redisOpt, broadcastChannel)

	// websocket handler
	hub := ws_worker.NewHub(&taskProcessor, &taskDistributor, broadcastChannel)
	wsHandler := ws_worker.NewHandler(hub)
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
