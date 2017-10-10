package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func Test_NewMetricsService_Should_Return_No_Error_When_OK(t *testing.T) {
	// Arrange
	s := &mockedService{}
	s.On("MakeMockedEndpoint").Return(nil).Once()
	h := NewMetricsService("test_ok", s.MakeMockedEndpoint())
	// Act
	err := h(context.Background(), &sarama.ConsumerMessage{Value: testEncodedMessage})
	// Assert
	assert.NoError(t, err)
}

func Test_NewMetricsService_Should_Return_Error_When_Internal_Error(t *testing.T) {
	// Arrange
	expected := errors.New("internal error")
	s := &mockedService{}
	s.On("MakeMockedEndpoint").Return(expected).Once()
	h := NewMetricsService("test_ko", s.MakeMockedEndpoint())
	// Act
	err := h(context.Background(), &sarama.ConsumerMessage{Value: testEncodedMessage})
	// Assert
	assert.Equal(t, expected, err)
}
