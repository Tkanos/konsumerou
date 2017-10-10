package kafka-consumer-router

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewServiceTracing_Should_Create_New_Service_Instance(t *testing.T) {
	//Arrange
	fake := new(mockedService)

	//Act
	tracer := NewServiceTracing(fake)

	//Assert
	assert.NotNil(t, tracer)
}

func Test_Tracing_Track_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessage{MyServiceID: id}
	fake.On("Track", context.Background(), model).Return(nil)
	tracer := NewServiceTracing(fake)

	//Act
	err := tracer.Track(context.Background(), model)

	//Assert
	assert.NoError(t, err)
}

func Test_Tracing_Track_With_Error_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessage{MyServiceID: id}
	fake.On("Track", context.Background(), model).Return(errors.New("test"))
	tracer := NewServiceTracing(fake)

	//Act
	err := tracer.Track(context.Background(), model)

	//Assert
	assert.Error(t, err)
}
