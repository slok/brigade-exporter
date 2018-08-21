package brigade

// Project is a representation of a brigade Project required by
// the application.
type Project struct {
	ID         string
	Name       string
	Repository string
	Namespace  string
	Worker     string
}
