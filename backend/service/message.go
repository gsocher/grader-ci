package service

import (
	"fmt"

	"github.com/dpolansky/grader-ci/amqp"
	"github.com/dpolansky/grader-ci/model"
	"github.com/sirupsen/logrus"

	"encoding/json"
)

type BuildMessageService interface {
	SendBuild(*model.BuildStatus) error
	ListenForBuildMessages()
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

func (m *messageService) ListenForBuildMessages() {
	m.messenger.ReadFromQueueWithCallback(model.AMQPStatusQueue, func(b []byte) {
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
	}, nil)
}

func (m *messageService) SendBuild(status *model.BuildStatus) error {
	status, err := m.buildService.UpdateBuild(status)
	if err != nil {
		return fmt.Errorf("Failed to update status for build %+v: %v", m, err)
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("Failed to marshal build to bytes: %v", err)
	}

	err = m.messenger.SendToQueue(model.AMQPBuildQueue, bytes)
	if err != nil {
		return fmt.Errorf("Failed to send build to queue: %v", err)
	}

	logrus.WithField("id", status.ID).Infof("Sent build")
	return nil
}
