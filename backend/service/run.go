package service

import (
	"fmt"

	"github.com/dpolansky/grader-ci/amqp"
	"github.com/dpolansky/grader-ci/model"
	"github.com/sirupsen/logrus"

	"encoding/json"
)

type BuildRunner interface {
	RunBuild(*model.BuildStatus) (*model.BuildStatus, error)
	ListenForUpdates()
}

type runner struct {
	m amqp.Messenger
	w BuildService
}

func NewAMQPBuildRunner(m amqp.Messenger, w BuildService) (BuildRunner, error) {
	return &runner{
		m: m,
		w: w,
	}, nil
}

func (r *runner) ListenForUpdates() {
	r.m.ReadFromQueueWithCallback(model.AMQPStatusQueue, func(b []byte) {
		var build model.BuildStatus
		err := json.Unmarshal(b, &build)
		if err != nil {
			logrus.WithError(err).WithField("body", string(b)).Errorf("Failed to unmarshal status update")
			return
		}

		_, err = r.w.UpdateBuild(&build)
		if err != nil {
			logrus.WithError(err).WithField("id", build.ID).Errorf("Failed to update build")
		}
	}, nil)
}

func (r *runner) RunBuild(m *model.BuildStatus) (*model.BuildStatus, error) {
	m, err := r.w.UpdateBuild(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to update status for build %+v: %v", m, err)
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal build to bytes: %v", err)
	}

	err = r.m.SendToQueue(model.AMQPBuildQueue, bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to send build to queue: %v", err)
	}

	logrus.WithField("id", m.ID).Infof("Sent build")
	return m, nil
}
