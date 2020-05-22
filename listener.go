package konsumerou

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

// Handler that handle kafka messages received
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Handlers defines a handler for a given topic
type Handlers map[string]Handler

// listener object represents kafka customer
type listener struct {
	consumer sarama.ConsumerGroup
	handlers Handlers
	ctx      context.Context
	topics   []string
}

const (
	contextTopicKey     = "topic"
	contextkeyKey       = "key"
	contextOffsetKey    = "offset"
	contextTimestampKey = "timestamp"
)

// Listener ...
type Listener interface {
	Subscribe() error
	Close()
}

// NewListener ...
func NewListener(ctx context.Context, brokers []string, groupID string, topicList string, handler Handler, config *sarama.Config) (Listener, error) {
	if brokers == nil || len(brokers) == 0 {
		return nil, errors.New("cannot create new listener, brokers cannot be empty")
	}
	if groupID == "" {
		return nil, errors.New("cannot create new listener, groupID cannot be empty")
	}
	if len(topicList) == 0 {
		return nil, errors.New("cannot create new listener, handlers cannot be empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	topics := strings.Split(topicList, ",")

	// create a map of handlers
	handlers := make(map[string]Handler)
	for _, topic := range topics {
		handlers[topic] = handler
	}

	// Init config
	if config == nil {
		config = sarama.NewConfig()
	}

	//config.Logger = logger //verbose mode
	config.Consumer.Return.Errors = true

	// Init consumer, consume errors & messages
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}
	consumer, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}
	//cluster.NewConsumer(brokers, groupID, topics, config)

	return &listener{
		consumer: consumer,
		handlers: handlers,
		ctx:      ctx,
		topics:   topics,
	}, nil
}

// Subscribe ...
func (l *listener) Subscribe() error {
	if l.consumer == nil {
		return errors.New("cannot subscribe. Customer is nil")
	}

	go func() {
		// When a session is over, make consumer join a new session, as long as the context is not cancelled
		for {
			// Consume make this consumer join the next session
			// This block until the `session` is over. (basically until next rebalance)
			err := l.consumer.Consume(l.ctx, l.topics, l)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			}
			if err := l.ctx.Err(); err != nil {
				// Check if context is cancelled
				return
			}
		}
	}()

	return nil
}

// ConsumerGroupHandler instances are used to handle individual topic/partition claims.
// It also provides hooks for your consumer group session life-cycle and allow you to
// trigger logic before or after the consume loop(s).
//
// PLEASE NOTE that handlers are likely be called from several goroutines concurrently,
// ensure that all state is safely protected against race conditions.

// Setup is run at the beginning of a new session, before ConsumeClaim
func (l *listener) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (l *listener) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (l *listener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		l.onNewMessage(msg, session)
	}
	return nil
}

func (l *listener) onNewMessage(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) {
	ctx := context.WithValue(l.ctx, contextTopicKey, msg.Topic)
	ctx = context.WithValue(ctx, contextkeyKey, msg.Key)
	ctx = context.WithValue(ctx, contextOffsetKey, msg.Offset)
	ctx = context.WithValue(ctx, contextTimestampKey, msg.Timestamp)

	err := l.handlers[msg.Topic](ctx, msg)

	if err != nil {
		// error should be handle on the handler
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}

	session.MarkMessage(msg, "")
}

func (l *listener) Close() {
	if l.consumer != nil {
		l.consumer.Close()
	}
}
