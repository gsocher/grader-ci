package service

import (
	"fmt"

	"github.com/dpolansky/grader-ci/pkg/amqp"
	"github.com/dpolansky/grader-ci/pkg/model"
	"github.com/sirupsen/logrus"

	"encoding/json"
)

type BuildMessageService interface {
	SendBuild(*model.BuildStatus) error
	ListenForBuildMessages(die chan struct{}) <-chan int
}

type messageService struct {
	messenger    amqp.Messenger
	buildService BuildService
}

func NewAMQPBuildMessageService(messenger amqp.Messenger, buildService BuildService) (BuildMessageService, error) {
	return &messageService{
		messenger:    messenger,
		buildService: buildService,
	}, nil
}

// returns a channel that can be listened to for IDs of received build updates
func (m *messageService) ListenForBuildMessages(die chan struct{}) <-chan int {
	received := make(chan int)

	go func(received chan int, die chan struct{}) {
		callback := func(b []byte) {
			var build model.BuildStatus

			err := json.Unmarshal(b, &build)
			if err != nil {
				logrus.WithError(err).WithField("body", string(b)).Errorf("Failed to unmarshal status update")
				return
			}

			_, err = m.buildService.UpdateBuild(&build)
			if err != nil {
				logrus.WithError(err).WithField("id", build.ID).Errorf("Failed to update build")
			}

			received <- build.ID
		}
		m.messenger.ReadFromQueueWithCallback(model.AMQPStatusQueue, callback, die)
	}(received, die)

	return received
}

func (m *messageService) SendBuild(status *model.BuildStatus) error {
	updated, err := m.buildService.UpdateBuild(status)
	if err != nil {
		return fmt.Errorf("Failed to update status for build %+v: %v", m, err)
	}

	bytes, err := json.Marshal(updated)
	if err != nil {
		return fmt.Errorf("Failed to marshal build to bytes: %v", err)
	}

	err = m.messenger.SendToQueue(model.AMQPBuildQueue, bytes)
	if err != nil {
		return fmt.Errorf("Failed to send build to queue: %v", err)
	}

	logrus.WithField("id", updated.ID).Infof("Sent build")
	return nil
}
