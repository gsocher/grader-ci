package amqp

import (
	"bytes"
	"testing"
)

const testingQueueName = "testing"

func TestSendReceive(t *testing.T) {
	c, err := NewAMQPDefaultClient()
	if err != nil {
		t.Fatalf("couldn't start client: %v", err)
	}

	expected := "foo"
	b := bytes.NewBufferString(expected)

	if err = c.SendToQueue(testingQueueName, b.Bytes()); err != nil {
		t.Fatalf("failed to send msg to queue: %v", err)
	}

	result := make(chan string)
	go c.ReadFromQueueWithCallback(testingQueueName, func(byt []byte) {
		result <- string(byt)
	}, nil)

	actual := <-result
	if actual != expected {
		t.Fatalf("expected %v got %v", expected, actual)
	}

	err = c.PurgeQueue(testingQueueName)
	if err != nil {
		t.Errorf("failed to purge queue: %v", err)
	}
}
