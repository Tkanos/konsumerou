//+build unit

package loginfailed

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testCorruptedMessage = []byte{1, 2, 3}
var testEncodedMessage = []byte(`{"my_service_id":"123"}`)

func Test_MakeMyServiceEndpoint_Should_Return_Object(t *testing.T) {
	// Arrange
	s := &mockedService{}
	// Act
	e := MakeMyServiceEndpoint(s)
	// Assert
	assert.NotNil(t, e)
}

func Test_MyServiceEndpoint_Should_Call_ProcessMessage(t *testing.T) {
	// Arrange
	s := &mockedService{}
	s.On("ProcessMessage", context.Background(), mock.AnythingOfType("MyServiceMessage")).Return(nil).Once()
	e := MakeMyServiceEndpoint(s)
	if assert.NotNil(t, e) {
		// Act
		err := e(&sarama.ConsumerMessage{Value: testEncodedMessage})
		// Assert
		assert.NoError(t, err)
		s.AssertExpectations(t)
	}
}

func Test_MyServiceEndpoint_Should_Return_Error_When_Corrupted_Message(t *testing.T) {
	// Arrange
	s := &mockedService{}
	e := MakeMyServiceEndpoint(s)
	if assert.NotNil(t, e) {
		// Act
		err := e(&sarama.ConsumerMessage{Value: testCorruptedMessage})
		// Assert
		assert.Error(t, err)
		s.AssertNumberOfCalls(t, "ProcessMessage", 0)
	}
}

func Test_DecodeMessage_Should_Return_Object(t *testing.T) {
	// Arrange
	expected := &MyServiceMessage{
		MyServiceID: "123",
	}
	// Act
	m, err := decodeMessage(testEncodedMessage)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
}

func Test_DecodeMessage_Should_Return_Error_When_Corrupted_Message(t *testing.T) {
	// Act

	m, err := decodeMessage(testCorruptedMessage)
	// Assert
	assert.Error(t, err)
	assert.Nil(t, m)
}
