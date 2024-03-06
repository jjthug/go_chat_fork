package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"server/db"
	"server/internal/user"
	"server/router"
	"server/ws"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	redisOpt := asynq.RedisClientOpt{
		Addr: "0.0.0.0:6379",
	}

	broadcastChannel := make(chan *ws.Message, 5)
	taskDistributor := ws.NewRedisTaskDistributor(redisOpt)

	if taskDistributor == nil {
		fmt.Errorf("taskDistributor is nil")
	}
	taskProcessor := ws.NewRedisTaskProcessor(redisOpt, broadcastChannel)

	// websocket handler
	hub := ws.NewHub(taskDistributor, broadcastChannel)
	wsHandler := ws.NewHandler(hub)
	go hub.Run()
	go taskProcessor.Start()

	//hub := ws.NewHub()
	//wsHandler := ws.NewHandler(hub)
	//go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
