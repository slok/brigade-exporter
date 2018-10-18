package brigade

import (
	"sync"
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
	builds, err := b.client.GetBuilds()
	if err != nil {
		return []*Job{}, err
	}

	// WARNING: N:M query.
	jobsC := make(chan *Job)
	var wg sync.WaitGroup
	wg.Add(len(builds))

	for _, bld := range builds {
		bld := bld
		go func() {
			defer wg.Done()
			b.getBuildJobs(bld, jobsC)
		}()
	}

	// Receive job results, the range will stop when closing
	// the channel (this is when all the goroutines have finished).
	var jobs []*Job
	go func() {
		for job := range jobsC {
			jobs = append(jobs, job)
		}
	}()

	// Wait until finished.
	wg.Wait()
	close(jobsC)

	return jobs, nil
}

func (b *brigade) getBuildJobs(build *azurebrigade.Build, jobsC chan<- *Job) {
	bjobs, err := b.client.GetBuildJobs(build)
	if err != nil {
		b.logger.Errorf("error retrieving job from build %s: %s", build.ID, err)
		return
	}

	for _, job := range bjobs {
		if job == nil {
			continue
		}

		jobsC <- &Job{
			ID:       job.ID,
			BuildID:  build.ID,
			Name:     job.Name,
			Image:    job.Image,
			Status:   job.Status.String(),
			Duration: b.getJobDuration(job),
			Creation: job.CreationTime,
			Start:    job.StartTime,
		}
	}
}
func (b *brigade) getJobDuration(job *azurebrigade.Job) time.Duration {
	if job == nil {
		return 0
	}

	var duration time.Duration

	// Only get duration if build finished.
	if job.Status == azurebrigade.JobSucceeded || job.Status == azurebrigade.JobFailed {
		duration = job.EndTime.Sub(job.StartTime)
	}

	// Only return if is a valid duration.
	if duration > 0 {
		return duration
	}

	return 0
}
