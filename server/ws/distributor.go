package ws

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendMessage(
		ctx context.Context,
		payload *PayloadSendMessage,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(residOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(residOpt)
	return &RedisTaskDistributor{client}
}
