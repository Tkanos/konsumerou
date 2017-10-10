package myservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewServiceLogging_Should_Create_New_Service_Instance(t *testing.T) {
	//Arrange
	fake := new(mockedService)

	//Act
	tracer := NewServiceLogging(fake)

	//Assert
	assert.NotNil(t, tracer)
}

func Test_Logging_ProcessMessage_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessage{MyServiceID: "123"}
	fake.On("ProcessMessage", context.Background(), &model).Return(nil)
	logger := NewServiceLogging(fake)

	//Act
	err := logger.ProcessMessage(context.Background(), &model)

	//Assert
	assert.NoError(t, err)
}

func Test_Logging_ProcessMessage_With_Error_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessage{MyServiceID: "123"}
	fake.On("ProcessMessage", context.Background(), &model).Return(errors.New("test"))
	logger := NewServiceLogging(fake)

	//Act
	err := logger.ProcessMessage(context.Background(), &model)

	//Assert
	assert.Error(t, err)
}
