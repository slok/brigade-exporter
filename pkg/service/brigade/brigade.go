package brigade

import (
	"github.com/Azure/brigade/pkg/storage"

	"github.com/slok/brigade-exporter/pkg/log"
)

// Interface is the interface that knows how to get data from brigade
// so the collectors can get the data.
type Interface interface {
	GetProjects() ([]*Project, error)
}

type brigade struct {
	client storage.Store
	logger log.Logger
}

// New returns a new brigade.Interface implementation.
func New(client storage.Store, logger log.Logger) Interface {
	return &brigade{
		client: client,
		logger: logger,
	}
}

func (b *brigade) GetProjects() ([]*Project, error) {
	bprs, err := b.client.GetProjects()
	if err != nil {
		return []*Project{}, err
	}

	prs := make([]*Project, len(bprs))
	for i, pr := range bprs {
		prs[i] = &Project{
			ID:         pr.ID,
			Name:       pr.Name,
			Repository: pr.Repo.Name,
			Namespace:  pr.Kubernetes.Namespace,
			Worker:     pr.Worker.Image(),
		}
	}

	return prs, nil
}
