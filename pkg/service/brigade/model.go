package brigade

import "time"

// Project is a representation of a brigade Project required by
// the application.
type Project struct {
	ID         string
	Name       string
	Repository string
	Namespace  string
	Worker     string
}

// Build is a representation of a brigade build required by the application.
type Build struct {
	ID        string
	ProjectID string
	Type      string
	Provider  string
	Version   string
	Status    string
	Duration  time.Duration
}

// Job is a representation of a brigade build job required by the application.
type Job struct {
	ID       string
	BuildID  string
	Name     string
	Image    string
	Status   string
	Duration time.Duration
	Creation time.Time
	Start    time.Time
}
