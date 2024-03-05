package ws_worker

import (
	"context"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
)

type TaskProcessor interface {
	ProcessTaskSendMessage(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server           *asynq.Server
	broadcastChannel chan *Message
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, broadcastChannel chan *Message) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{})
	//Queues: map[string]int{QueueCritical: 10},

	return &RedisTaskProcessor{
		server:           server,
		broadcastChannel: broadcastChannel,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendMessage, processor.ProcessTaskSendMessage)

	return processor.server.Start(mux)
}
