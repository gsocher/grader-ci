package repo

import (
	"fmt"
	"sort"
	"sync"

	"github.com/dpolansky/ci/model"
)

type BuildStatusRepo interface {
	UpsertStatus(build *model.BuildStatus) (*model.BuildStatus, error)
	GetStatusForId(id int) (*model.BuildStatus, error)
	GetStatuses() ([]*model.BuildStatus, error)
}

type memStatusRepo struct {
	lock   *sync.Mutex
	builds []*model.BuildStatus
}

func NewInMemoryStatusRepo() BuildStatusRepo {
	return &memStatusRepo{
		lock:   &sync.Mutex{},
		builds: []*model.BuildStatus{},
	}
}

func (m *memStatusRepo) UpsertStatus(build *model.BuildStatus) (*model.BuildStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for i, b := range m.builds {
		fmt.Printf("%v : %#v\n", i, b)
	}

	// if we haven't seen this build yet, generate a new id
	if build.ID == 0 {
		build.ID = m.getNextStatusID()
		m.builds = append(m.builds, build)
		return build, nil
	}

	m.builds[build.ID-1] = build
	return build, nil
}

func (m *memStatusRepo) GetStatusForId(id int) (*model.BuildStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if id <= 0 || id > len(m.builds) {
		return nil, fmt.Errorf("No build status with id %v exists", id)
	}

	return m.builds[id-1], nil
}

func (m *memStatusRepo) GetStatuses() ([]*model.BuildStatus, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var buildsTyped Builds = m.builds
	var buildsCopy Builds = make([]*model.BuildStatus, len(buildsTyped))
	copy(buildsCopy, buildsTyped)
	sort.Sort(buildsCopy)

	return buildsCopy, nil
}

func (m *memStatusRepo) getNextStatusID() int {
	return len(m.builds) + 1
}

type Builds []*model.BuildStatus

func (slice Builds) Len() int {
	return len(slice)
}

func (slice Builds) Less(i, j int) bool {
	return slice[i].ID > slice[j].ID
}

func (slice Builds) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
