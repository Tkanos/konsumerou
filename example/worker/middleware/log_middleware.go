package middleware

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/tkanos/konsumerou"
)

// Logger represents some logger with methods to log info
type Logger interface {
	Printf(format string, v ...interface{})
}

type logService struct {
	logger Logger
}

// NewLogService creates a layer of service that add logging capability
func NewLogService(l Logger, next konsumerou.Handler) konsumerou.Handler {
	return logMiddleware(l).logging(next)
}

func logMiddleware(l Logger) *logService {
	return &logService{l}
}

func (s logService) logging(next konsumerou.Handler) konsumerou.Handler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (err error) {
		if msg.Value == nil {
			err = errors.New("cannot process empty message")
		} else {
			err = next(ctx, msg)
		}
		if err != nil {
			s.logger.Printf("cannot process offset message %v, error : %v", msg.Offset, err)
		}
		return
	}
}
