package konsumerou

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewListener_EmptyBrocker(t *testing.T) {
	listener, err := NewListener(context.Background(), []string{}, "group-id", "my-topic", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}

func Test_NewListener_EmptyGroup(t *testing.T) {
	listener, err := NewListener(context.Background(), []string{"brocker1", "brocker2"}, "", "my-topic", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}

func Test_NewListener_EmptyTopicList(t *testing.T) {
	listener, err := NewListener(context.Background(), []string{"brocker1", "brocker2"}, "group-id", "", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}

func Test_NewListenerHandlers_EmptyHandlers(t *testing.T) {
	l, err := NewListenerHandlers(context.Background(), []string{"brocker1", "brocker2"}, "groupID", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, l)
}
