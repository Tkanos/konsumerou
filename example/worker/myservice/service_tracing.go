package myservice

/*
import (
	"context"

	tracing "github.com/ricardo-ch/go-tracing"
)

// MyServiceTracing ...
type MyServiceTracing struct {
	next MyServiceMessageProcessor
}

// NewServiceTracing ...
func NewServiceTracing(l MyServiceMessageProcessor) MyServiceMessageProcessor {
	return MyServiceTracing{
		next: l,
	}
}

// ProcessMessage ...
func (s MyServiceTracing) ProcessMessage(ctx context.Context, msg *MyServiceMessage) (err error) {
	span, ctx := tracing.CreateSpan(ctx, "myservice.MyService::ProcessMessage", &map[string]interface{}{"ID": msg.MyServiceID})
	defer func() {
		if err != nil {
			tracing.SetSpanError(span, err)
		}
		span.Finish()
	}()

	return s.next.ProcessMessage(ctx, msg)
}
*/
