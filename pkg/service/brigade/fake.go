package brigade

import (
	"fmt"
	"time"

	azurebrigade "github.com/Azure/brigade/pkg/brigade"
)

var (
	fakedBuildEventTypes = []string{"push", "pull_request", "deploy", "deploy_post_hook", "tag", "debug"}
	fakedBuildProviders  = []string{"github", "docker", "gitlab", "brig", "toilet"}
	fakedJobStatus       = []azurebrigade.JobStatus{azurebrigade.JobPending, azurebrigade.JobRunning, azurebrigade.JobSucceeded, azurebrigade.JobFailed, azurebrigade.JobUnknown}
)

type fake struct{}

// NewFake returns a new fake implementation of the Brigade interface.
func NewFake() Interface {
	return &fake{}
}

func (f *fake) GetProjects() ([]*Project, error) {
	var prs []*Project

	for i := 0; i < 10; i++ {
		prs = append(prs, &Project{
			Name:       fmt.Sprintf("project-%d", i),
			ID:         fmt.Sprintf("prj-id-%d", i),
			Repository: fmt.Sprintf("github.com/fake-exporter/project-%d", i),
			Namespace:  fmt.Sprintf("ns%d", i),
			Worker:     fmt.Sprintf("brigade-worker-%d", i),
		})
	}
	return prs, nil
}

func (f *fake) GetBuilds() ([]*Build, error) {
	var blds []*Build

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			blds = append(blds, &Build{
				ID:        fmt.Sprintf("build-id-%d%d", i, j),
				ProjectID: fmt.Sprintf("prj-id-%d", i),
				Type:      fakedBuildEventTypes[time.Now().UnixNano()%int64(len(fakedBuildEventTypes))],
				Provider:  fakedBuildProviders[time.Now().UnixNano()%int64(len(fakedBuildProviders))],
				Version:   fmt.Sprintf("%d", time.Now().UnixNano()),
				Status:    fakedJobStatus[time.Now().UnixNano()%int64(len(fakedJobStatus))].String(),
				Duration:  time.Duration(time.Now().UnixNano()%4000) * time.Second,
			})

		}
	}
	return blds, nil
}

func (f *fake) GetJobs() ([]*Job, error) {
	var jobs []*Job

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			jobs = append(jobs, &Job{
				ID:       fmt.Sprintf("job-id-%d%d%d", i, j, i),
				Name:     fmt.Sprintf("job-%d%d%d", i, j, i),
				BuildID:  fmt.Sprintf("build-id-%d%d", i, j),
				Image:    fmt.Sprintf("fake/job-image:%d%d", i, j),
				Status:   fakedJobStatus[time.Now().UnixNano()%int64(len(fakedJobStatus))].String(),
				Duration: time.Duration(time.Now().UnixNano()%4000) * time.Second,
			})
		}
	}
	return jobs, nil
}
