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

	// With this ID we make it change every 10m
	startID := int(time.Now().Unix() / 600)
	// Change Status every 120m.
	statusSalt := (time.Now().Unix() / 120)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fakeIdentity := i + j
			statusRand := statusSalt * int64(j*i)

			blds = append(blds, &Build{
				ID:        fmt.Sprintf("build-id-%d%d%d", startID, i, j),
				ProjectID: fmt.Sprintf("prj-id-%d", i),
				Type:      fakedBuildEventTypes[(22*startID*fakeIdentity)%len(fakedBuildEventTypes)],
				Provider:  fakedBuildProviders[(23*startID*fakeIdentity)%len(fakedBuildProviders)],
				Version:   fmt.Sprintf("%d", (1234567 * startID * fakeIdentity)),
				Status:    fakedJobStatus[statusRand%int64(len(fakedJobStatus))].String(),
				Duration:  time.Duration((startID*fakeIdentity)%4000) * time.Second,
			})

		}
	}
	return blds, nil
}

func (f *fake) GetJobs() ([]*Job, error) {
	var jobs []*Job

	// With this ID we make it change every 10m
	startID := time.Now().Unix() / 600
	// Change Status every 120m.
	statusSalt := (time.Now().Unix() / 120)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 15; k++ {
				statusRand := statusSalt * int64(j*i*k)
				jobs = append(jobs, &Job{
					ID:       fmt.Sprintf("job-id-%d%d%d%d", startID, j, i, k),
					Name:     fmt.Sprintf("job-%d%d%d%d", startID, j, i, k),
					BuildID:  fmt.Sprintf("build-id-%d%d%d", startID, i, j),
					Image:    fmt.Sprintf("fake/job-image:%d%d", i, j),
					Status:   fakedJobStatus[statusRand%int64(len(fakedJobStatus))].String(),
					Duration: time.Duration((987654321*startID)%4000) * time.Second,
				})
			}
		}
	}
	return jobs, nil
}
