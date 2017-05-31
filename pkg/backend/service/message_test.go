package service

import (
	"encoding/json"
	"testing"

	"github.com/dpolansky/grader-ci/pkg/amqp"
	"github.com/dpolansky/grader-ci/pkg/backend/service/fakes"
	"github.com/dpolansky/grader-ci/pkg/model"
)

func TestListenForBuildMessages(t *testing.T) {
	fakeBuildService := new(fakes.FakeBuildService)

	ourClient, err := amqp.NewAMQPDefaultClient()
	if err != nil {
		t.Fatalf("failed to create amqp client: %v", err)
	}
	defer ourClient.PurgeQueue(model.AMQPStatusQueue)

	theirClient, err := amqp.NewAMQPDefaultClient()
	if err != nil {
		t.Fatalf("failed to create amqp client: %v", err)
	}

	msgService, err := NewAMQPBuildMessageService(ourClient, fakeBuildService)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// if they send a build update, we should receive it
	die := make(chan struct{})
	received := msgService.ListenForBuildMessages(die)

	status := &model.BuildStatus{ID: 0}
	b, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if err := theirClient.SendToQueue(model.AMQPStatusQueue, b); err != nil {
		t.Fatalf("their client failed to send: %v", err)
	}

	t.Logf("waiting for status on queue")
	id := <-received
	if id != status.ID {
		t.Fatalf("expected id %v got %v", status.ID, id)
	}

	die <- struct{}{}
}

func TestSendBuildMessage(t *testing.T) {
	fakeBuildService := new(fakes.FakeBuildService)

	ourClient, err := amqp.NewAMQPDefaultClient()
	if err != nil {
		t.Fatalf("failed to create amqp client: %v", err)
	}
	theirClient, err := amqp.NewAMQPDefaultClient()
	if err != nil {
		t.Fatalf("failed to create amqp client: %v", err)
	}
	defer ourClient.PurgeQueue(model.AMQPBuildQueue)

	msgService, err := NewAMQPBuildMessageService(ourClient, fakeBuildService)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// if we send a build message, they should receive it
	status := &model.BuildStatus{ID: 0}
	fakeBuildService.UpdateBuildReturns(status, nil)
	if err := msgService.SendBuild(status); err != nil {
		t.Fatalf("failed to send build: %v", err)
	}

	received := make(chan int)
	die := make(chan struct{})
	go func(received chan int, die chan struct{}) {
		callback := func(b []byte) {
			var s model.BuildStatus
			if err := json.Unmarshal(b, &s); err != nil {
				t.Fatalf("failed to unmarshal status: %v", err)
			}

			received <- s.ID
		}
		theirClient.ReadFromQueueWithCallback(model.AMQPBuildQueue, callback, die)
	}(received, die)

	id := <-received

	if id != status.ID {
		t.Fatalf("expected id %v got %v", status.ID, id)
	}

	die <- struct{}{}
}
