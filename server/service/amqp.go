package service

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Messenger can read from and write to AMQP queues.
type Messenger interface {
	ReadFromQueueWithCallback(queueName string, callback func(amqp.Delivery), die chan struct{}) error
	SendToQueue(queueName string, b []byte) error
}

// NewAMQPClient creates a new AMQP client and creates a connection with the given URL
func NewAMQPClient(url string) (Messenger, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial amqp: %v", err)
	}

	c := &amqpClient{
		conn: conn,
	}

	return c, nil
}

type amqpClient struct {
	conn *amqp.Connection
}

// SendToQueue adds a message to a queue.
func (c *amqpClient) SendToQueue(queueName string, b []byte) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %v", err)
	}

	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})

	if err != nil {
		return fmt.Errorf("Failed to publish: %v", err)
	}

	return nil
}

// ReadFromQueueWithCallback is a blocking call that reads messages from a queue and invokes a given callback function
// on each message until signaled to stop from the die channel.
func (c *amqpClient) ReadFromQueueWithCallback(queueName string, callback func(amqp.Delivery), die chan struct{}) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %v", err)
	}

	// Read messages from the channel until we are signaled to stop via die channel
	for {
		select {
		case m := <-msgs:
			callback(m)
		case <-die:
			return nil
		}
	}
}

func NewMockClient() Messenger {
	return &mockClient{}
}

type mockClient struct {
}

func (m *mockClient) SendToQueue(queueName string, b []byte) error {
	return nil
}

func (m *mockClient) ReadFromQueueWithCallback(queueName string, callback func(amqp.Delivery), die chan struct{}) error {
	return nil
}
