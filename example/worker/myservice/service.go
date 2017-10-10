package myservice

import "context"

type MyServiceMessageProcessor interface {
	ProcessMessage(context.Context, *MyServiceMessage) error
}

type service struct {
}

// NewService ....
func NewService() MyServiceMessageProcessor {
	return service{}
}

// ProcessMessage ...
func (s service) ProcessMessage(ctx context.Context, msg *MyServiceMessage) error {
	// PUT ALL YOUR LOGIC HERE
	return nil
}
