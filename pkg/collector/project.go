package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	projectSubSystem = "project"
)

// project is the Brigade project subcollector. this colletor will collect
// the metrics regarding brigade projects.
// Satisfies internfal collector interface.
type project struct {
	brigadeSVC brigade.Interface
	logger     log.Logger

	// Metrics.
	projectInfoDesc *prometheus.Desc
}

// NewProject returns a new project subcollector.
func NewProject(brigadeSVC brigade.Interface, logger log.Logger) subcollector {
	return &project{
		brigadeSVC: brigadeSVC,
		logger:     logger,

		projectInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, projectSubSystem, "info"),
			"Brigade project information.",
			[]string{"id", "name", "repository", "namespace", "worker"}, nil,
		),
	}
}

// Collect satisfies subcollector.
func (p *project) Collect(ch chan<- prometheus.Metric) error {
	// Collect project info.
	prs, err := p.brigadeSVC.GetProjects()
	if err != nil {
		return err
	}

	for _, pr := range prs {
		ch <- prometheus.MustNewConstMetric(
			p.projectInfoDesc,
			prometheus.GaugeValue,
			1,
			pr.ID, pr.Name, pr.Repository, pr.Namespace, pr.Worker)
	}

	return nil
}
