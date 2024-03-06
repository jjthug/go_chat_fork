package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
)

const TaskSendMessage = "task:send_message"

type PayloadSendMessage struct {
	Message Message `json:"message"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendMessage(ctx context.Context,
	payload *PayloadSendMessage,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}
	task := asynq.NewTask(TaskSendMessage, jsonPayload, opts...)
	if distributor == nil {
		// Handle the error, e.g., log it or initialize the client properly
		return errors.New("distributor is nil")
	}

	if distributor.client == nil {
		return errors.New("distributor.client is nil")
	}

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}

	log.Println("info=>", info)
	log.Println("payload=>", string(task.Payload()))
	//log.Info().Str("type", task.Type())
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendMessage(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendMessage
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload %v", asynq.SkipRetry)
	}

	// Send message to hub channel TODO check
	processor.broadcastChannel <- &payload.Message

	return nil
}
