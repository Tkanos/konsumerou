package konsumerou

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

// Handler that handle kafka messages received
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Handlers defines a handler for a given topic
type Handlers map[string]Handler

// listener object represents kafka customer
type listener struct {
	consumer *cluster.Consumer
	handlers Handlers
}

// Listener ...
type Listener interface {
	Subscribe(exit chan bool) error
	Close()
}

func NewListener(brokers []string, groupID string, topicList string, handler Handler, config *cluster.Config) (Listener, error) {
	// create a map of handlers
	handlers := make(map[string]Handler)
	for _, topic := range strings.Split(topicList, ",") {
		handlers[topic] = handler
	}

	return NewListenerHandlers(brokers, groupID, handlers, config)
}

func NewListenerHandlers(brokers []string, groupID string, handlers Handlers, config *cluster.Config) (Listener, error) {
	if brokers == nil || len(brokers) == 0 {
		return nil, errors.New("cannot create new listener, brokers cannot be empty")
	}
	if groupID == "" {
		return nil, errors.New("cannot create new listener, groupID cannot be empty")
	}
	if handlers == nil {
		return nil, errors.New("cannot create new listener, handlers cannot be empty")
	}

	// Init config
	if config == nil {
		config = cluster.NewConfig()
	}

	//config.Logger = logger //verbose mode
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// Init consumer, consume errors & messages
	var topics []string
	for k := range handlers {
		topics = append(topics, k)
	}
	consumer, err := cluster.NewConsumer(brokers, groupID, topics, config)
	if err != nil {
		return nil, err
	}

	return &listener{
		consumer: consumer,
		handlers: handlers,
	}, nil
}

func (l *listener) Subscribe(exit chan bool) error {
	if l.consumer == nil {
		return errors.New("cannot subscribe. Customer is nil")
	}

	go func() {
		// Consume all channels, wait for signal to exit
		for {
			select {
			case msg, more := <-l.consumer.Messages():
				if more {
					if l.handlers[msg.Topic](context.Background(), msg) == nil {
						l.consumer.MarkOffset(msg, "")
					}
				}
			case ntf, more := <-l.consumer.Notifications():
				if more {
					fmt.Fprintf(os.Stdout, "Rebalanced: %+v\n", ntf)
				}
			case err, more := <-l.consumer.Errors():
				if more {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				}
			case <-exit:
				return
			}
		}
	}()

	return nil
}

func (l *listener) Close() {
	if l.consumer != nil {
		l.consumer.Close()
	}
}
