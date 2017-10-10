package myservice

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/tkanos/konsumerou"
)

// MakeMyServiceEndpoint ...
func MakeMyServiceEndpoint(s MyServiceMessageProcessor) konsumerou.Handler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		message, err := decodeMessage(msg.Value)
		if err != nil {
			return err
		}
		return s.ProcessMessage(ctx, message)
	}
}

func decodeMessage(msg []byte) (*MyServiceMessage, error) {

	message := MyServiceMessage{}
	if err := json.Unmarshal(msg, &message); err != nil {
		return nil, err
	}
	return &message, nil
}
