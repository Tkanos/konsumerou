package myservice

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockedService struct {
	mock.Mock
}

func (m *mockedService) ProcessMessage(ctx context.Context, msg *MyServiceMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
