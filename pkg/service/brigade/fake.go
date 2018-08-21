package brigade

import "fmt"

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
			ID:         fmt.Sprintf("id-%d", i),
			Repository: fmt.Sprintf("github.com/fake-exporter/project-%d", i),
			Namespace:  fmt.Sprintf("ns%d", i),
			Worker:     fmt.Sprintf("brigade-worker-%d", i),
		})
	}
	return prs, nil
}
