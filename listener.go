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

//Handler that handle kafka messages received
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// listener object represents kafka customer
type listener struct {
	consumer *cluster.Consumer
}

// Listener ...
type Listener interface {
	Subscribe(handler Handler, exit chan bool) error
	Close()
}

// NewListener ...
func NewListener(brokers []string, groupID string, topicList string, offset int64, config *cluster.Config) (Listener, error) {
	if brokers == nil || len(brokers) == 0 {
		return nil, errors.New("cannot create new listener, brockers cannot be empty")
	}
	if groupID == "" {
		return nil, errors.New("cannot create new listener, groupID cannot be empty")
	}
	if topicList == "" {
		return nil, errors.New("cannot create new listener, topicList cannot be empty")
	}

	// Init config
	if config == nil {
		config = cluster.NewConfig()
	}

	//config.Logger = logger //verbose mode
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = offset

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(brokers, groupID, strings.Split(topicList, ","), config)
	if err != nil {
		return nil, err
	}

	return &listener{
		consumer: consumer,
	}, nil
}

func (l *listener) Subscribe(handler Handler, exit chan bool) error {
	if l.consumer == nil {
		return errors.New("cannot subscribe. Customer is nil")
	}

	go func() {
		// Consume all channels, wait for signal to exit
		for {
			select {
			case msg, more := <-l.consumer.Messages():
				if more {
					if handler(context.Background(), msg) == nil {
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
