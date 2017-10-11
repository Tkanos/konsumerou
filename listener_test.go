package konsumerou

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewListener_EmptyBrocker(t *testing.T) {
	listener, err := NewListener([]string{}, "my-topic", "group-id", 0, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}

func Test_NewListener_EmptyGroup(t *testing.T) {
	listener, err := NewListener([]string{"brocker1", "brocker2"}, "", "my-topic", 0, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}

func Test_NewListener_EmptyTopicList(t *testing.T) {
	listener, err := NewListener([]string{"brocker1", "brocker2"}, "group-id", "", 0, nil)
	assert.Error(t, err)
	assert.Nil(t, listener)
}
