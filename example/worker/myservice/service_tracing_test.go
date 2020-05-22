package myservice

/*
import (
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

func Test_Tracing_ProcessMessage_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessageProcessor{MyServiceID: "123"}
	fake.On("ProcessMessage", mock.Anything, &model).Return(nil)
	tracer := NewServiceTracing(fake)

	//Act
	err := tracer.ProcessMessage(context.Background(), &model)

	//Assert
	assert.NoError(t, err)
}

func Test_Tracing_ProcessMessage_With_Error_Should_Go_Throught_The_Method(t *testing.T) {
	//Arrange
	fake := new(mockedService)
	model := MyServiceMessageProcessor{MyServiceID: "123"}
	fake.On("ProcessMessage", mock.Anything, &model).Return(errors.New("test"))
	tracer := NewServiceTracing(fake)

	//Act
	err := tracer.ProcessMessage(context.Background(), &model)

	//Assert
	assert.Error(t, err)
}
*/
