package myservice

import (
	"context"
	"fmt"
	"os"

	"github.com/tkanos/generator-consumer-kafka-go/config"
)

// myServiceLogger ...
type myServiceLogger struct {
	next MyServiceMessageProcessor
}

// NewServiceLogging ...
func NewServiceLogging(s MyServiceMessageProcessor) MyServiceMessageProcessor {
	return myServiceLogger{
		next: s,
	}
}

// ProcessMessage ...
func (s myServiceLogger) ProcessMessage(ctx context.Context, msg *MyServiceMessage) error {

	if config.Config.Verbose {
		if msg != nil {
			fmt.Fprintf(os.Stdout, "%s\n", msg.MyServiceID)
		}
	}

	return s.next.ProcessMessage(ctx, msg)
}
