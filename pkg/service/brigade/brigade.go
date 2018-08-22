package brigade

import (
	"time"

	azurebrigade "github.com/Azure/brigade/pkg/brigade"
	"github.com/Azure/brigade/pkg/storage"

	"github.com/slok/brigade-exporter/pkg/log"
)

// Interface is the interface that knows how to get data from brigade
// so the collectors can get the data.
type Interface interface {
	GetProjects() ([]*Project, error)
	GetBuilds() ([]*Build, error)
	GetJobs() ([]*Job, error)
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

		// Only set image if there is an image name.
		image := ""
		if pr.Worker.Name != "" {
			image = pr.Worker.Image()
		}

		prs[i] = &Project{
			ID:         pr.ID,
			Name:       pr.Name,
			Repository: pr.Repo.Name,
			Namespace:  pr.Kubernetes.Namespace,
			Worker:     image,
		}
	}

	return prs, nil
}

func (b *brigade) GetBuilds() ([]*Build, error) {
	bblds, err := b.client.GetBuilds()
	if err != nil {
		return []*Build{}, err
	}

	blds := make([]*Build, len(bblds))
	for i, bld := range bblds {
		blds[i] = &Build{
			ID:        bld.ID,
			ProjectID: bld.ProjectID,
			Type:      bld.Type,
			Provider:  bld.Provider,
			Version:   bld.Revision.Commit,
			Status:    b.getBuildStatus(bld),
			Duration:  b.getBuildDuration(bld),
		}
	}

	return blds, nil
}

func (b *brigade) getBuildStatus(bld *azurebrigade.Build) string {
	if bld.Worker == nil {
		return azurebrigade.JobUnknown.String()
	}

	return bld.Worker.Status.String()
}

func (b *brigade) getBuildDuration(bld *azurebrigade.Build) time.Duration {
	if bld.Worker == nil {
		return 0
	}

	var duration time.Duration

	// Only get duration if build finished.
	if bld.Worker.Status == azurebrigade.JobSucceeded || bld.Worker.Status == azurebrigade.JobFailed {
		duration = bld.Worker.EndTime.Sub(bld.Worker.StartTime)
	}

	// Only return if is a valid duration.
	if duration > 0 {
		return duration
	}

	return 0
}

func (b *brigade) GetJobs() ([]*Job, error) {
	return []*Job{}, nil
}
