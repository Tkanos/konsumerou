package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tkanos/konsumerou"
)

// Represents an encoded message
var testEncodedMessage = []byte{10, 3, 49, 50, 51}

type mockedService struct {
	mock.Mock
}

func (s *mockedService) MakeMockedEndpoint() konsumerou.Handler {
	args := s.Called()
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (err error) {
		return args.Error(0)
	}
}

type mockedLogger struct {
	mock.Mock
}

func (l *mockedLogger) Printf(format string, v ...interface{}) {
	l.Called()
	return
}

func Test_NewLogService_Should_Return_No_Error_When_OK(t *testing.T) {
	// Arrange
	s := &mockedService{}
	s.On("MakeMockedEndpoint").Return(nil).Once()
	l := &mockedLogger{}
	h := NewLogService(l, s.MakeMockedEndpoint())
	// Act
	err := h(context.Background(), &sarama.ConsumerMessage{Value: testEncodedMessage})
	// Assert
	assert.NoError(t, err)
	l.AssertNumberOfCalls(t, "Printf", 0)
}

func Test_NewLogService_Should_Log_Error_When_Message_Is_Empty(t *testing.T) {
	// Arrange
	s := &mockedService{}
	s.On("MakeMockedEndpoint").Return(nil).Once()
	l := &mockedLogger{}
	l.On("Printf").Return(nil).Once()
	h := NewLogService(l, s.MakeMockedEndpoint())
	// Act
	err := h(context.Background(), &sarama.ConsumerMessage{Value: nil})
	// Assert
	assert.Error(t, err)
	l.AssertNumberOfCalls(t, "Printf", 1)
}

func Test_NewLogService_Should_Log_Error_When_Internal_Error(t *testing.T) {
	// Arrange
	expected := errors.New("internal error")
	s := &mockedService{}
	s.On("MakeMockedEndpoint").Return(expected).Once()
	l := &mockedLogger{}
	l.On("Printf").Return(nil).Once()
	h := NewLogService(l, s.MakeMockedEndpoint())
	// Act
	err := h(context.Background(), &sarama.ConsumerMessage{Value: testEncodedMessage})
	// Assert
	assert.Equal(t, expected, err)
	l.AssertNumberOfCalls(t, "Printf", 1)
}
