package server

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"sync"

	"github.com/dpolansky/ci/model"
	"github.com/dpolansky/ci/server/amqp"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	amqpAPI "github.com/streadway/amqp"
)

// BuildService represents the logic for starting and checking the status of builds.
type BuildService interface {
	StartBuild(cloneURL string) (*model.BuildStatus, error)
	GetStatusForBuild(id string) (*model.BuildStatus, error)
	UpdateStatusForBuild(build *model.BuildStatus) *model.BuildStatus
	ListenForUpdates()
}

type buildService struct {
	log        *logrus.Entry
	amqpClient amqp.ReadWriter

	lock   *sync.Mutex
	builds map[string]*model.BuildStatus
}

func NewBuildService(amqpClient amqp.ReadWriter) BuildService {
	log := logrus.WithField("module", "BuildService")

	return &buildService{
		log:        log,
		amqpClient: amqpClient,
		lock:       &sync.Mutex{},
		builds:     map[string]*model.BuildStatus{},
	}
}

func (b *buildService) ListenForUpdates() {
	b.log.Infof("Listening for status updates")

	b.amqpClient.ReadFromQueueWithCallback(model.AMQPStatusQueue, func(m amqpAPI.Delivery) {
		var build model.BuildStatus
		err := json.Unmarshal(m.Body, &build)
		if err != nil {
			b.log.WithError(err).WithField("body", string(m.Body)).Errorf("Failed to unmarshal status update")
			return
		}

		b.UpdateStatusForBuild(&build)
	}, nil)
}

func (b *buildService) StartBuild(cloneURL string) (*model.BuildStatus, error) {
	id := uuid.New().String()

	build := &model.BuildStatus{
		ID:         id,
		CloneURL:   cloneURL,
		Language:   "golang", // TODO: eventually clone into the URL and parse a config file to read the language
		LastUpdate: time.Now(),
		Status:     model.StatusBuildWaiting,
	}

	bytes, err := json.Marshal(build)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal build to bytes: %v", err)
	}

	err = b.amqpClient.SendToQueue(model.AMQPBuildQueue, bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to send build to queue: %v", err)
	}

	b.log.WithField("id", id).Infof("Sent build")

	b.UpdateStatusForBuild(build)
	return build, nil
}

func (b *buildService) GetStatusForBuild(id string) (*model.BuildStatus, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	status, found := b.builds[id]
	if !found {
		return nil, fmt.Errorf("No build found with id %v", id)
	}

	return status, nil
}

func (b *buildService) UpdateStatusForBuild(build *model.BuildStatus) *model.BuildStatus {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.builds[build.ID] = build

	b.log.WithFields(logrus.Fields{
		"id":     build.ID,
		"status": build.Status,
	}).Infof("Build status updated")

	return build
}

func (b *buildService) GetAllBuilds() []*model.BuildStatus {
	b.lock.Lock()
	defer b.lock.Unlock()

	builds := Builds{}

	for _, build := range b.builds {
		builds = append(builds, build)
	}

	sort.Sort(builds)
	return builds
}

type Builds []*model.BuildStatus

func (slice Builds) Len() int {
	return len(slice)
}

func (slice Builds) Less(i, j int) bool {
	return slice[i].ID < slice[j].ID
}

func (slice Builds) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
