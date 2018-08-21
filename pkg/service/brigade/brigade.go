package brigade

// Interface is the interface that knows how to get data from brigade
// so the collectors can get the data.
type Interface interface {
	GetProjects() ([]*Project, error)
}

type brigade struct {
}

// New returns a new brigade.Interface implementation.
func New() Interface {
	return &brigade{}
}

func (b *brigade) GetProjects() ([]*Project, error) {
	return []*Project{}, nil
}
