package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/repo"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Builder represents the logic for starting and checking the status of builds.
type Builder interface {
	StartBuild(cloneURL, branch string) (*model.BuildStatus, error)
	GetStatusForBuild(id int) (*model.BuildStatus, error)
	UpdateStatusForBuild(build *model.BuildStatus) (*model.BuildStatus, error)
	GetBuildStatuses() ([]*model.BuildStatus, error)
	ListenForUpdates()
}

type buildService struct {
	log       *logrus.Entry
	messenger Messenger
	repo      repo.BuildStatusRepo
}

func NewBuilder(m Messenger, r repo.BuildStatusRepo) Builder {
	log := logrus.WithField("module", "BuildService")

	return &buildService{
		log:       log,
		messenger: m,
		repo:      r,
	}
}

func (b *buildService) ListenForUpdates() {
	b.log.Infof("Listening for status updates")

	b.messenger.ReadFromQueueWithCallback(model.AMQPStatusQueue, func(msg amqp.Delivery) {
		var build model.BuildStatus
		err := json.Unmarshal(msg.Body, &build)
		if err != nil {
			b.log.WithError(err).WithField("body", string(msg.Body)).Errorf("Failed to unmarshal status update")
			return
		}

		b.UpdateStatusForBuild(&build)
	}, nil)
}

func (b *buildService) StartBuild(cloneURL, branch string) (*model.BuildStatus, error) {
	build := &model.BuildStatus{
		CloneURL: cloneURL,
		Branch:   branch,
		Status:   model.StatusBuildWaiting,
	}

	status, err := b.UpdateStatusForBuild(build)
	if err != nil {
		return nil, fmt.Errorf("Failed to update status for build %+v: %v", build, err)
	}

	bytes, err := json.Marshal(build)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal build to bytes: %v", err)
	}

	err = b.messenger.SendToQueue(model.AMQPBuildQueue, bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to send build to queue: %v", err)
	}

	b.log.WithField("id", status.ID).Infof("Sent build")
	return build, nil
}

func (b *buildService) GetStatusForBuild(id int) (*model.BuildStatus, error) {
	return b.repo.GetStatusForId(id)
}

func (b *buildService) UpdateStatusForBuild(build *model.BuildStatus) (*model.BuildStatus, error) {
	build.LastUpdate = time.Now()
	build, err := b.repo.UpsertStatus(build)
	if err != nil {
		return nil, err
	}

	b.log.WithFields(logrus.Fields{
		"id":     build.ID,
		"status": build.Status,
	}).Infof("Build status updated")

	return build, nil
}

func (b *buildService) GetBuildStatuses() ([]*model.BuildStatus, error) {
	return b.repo.GetStatuses()
}
