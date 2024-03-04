package ws_worker

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	ProcessTaskSendMessage(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	hub    *Hub
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, hub *Hub) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{})

	return &RedisTaskProcessor{
		server: server,
		hub:    hub,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendMessage, processor.ProcessTaskSendMessage)

	return processor.server.Start(mux)
}
