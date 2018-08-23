package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	jobSubSystem = "job"
)

// job is the Brigade Job subcollector. this colletor will collect
// the metrics regarding brigade jobs.
// Satisfies internfal collector interface.
type job struct {
	brigadeSVC brigade.Interface
	logger     log.Logger

	// Metrics.
	jobInfoDesc     *prometheus.Desc
	jobStatusDesc   *prometheus.Desc
	jobDurationDesc *prometheus.Desc
}

// NewJob returns a new job subcollector.
func NewJob(brigadeSVC brigade.Interface, logger log.Logger) subcollector {
	return &job{
		brigadeSVC: brigadeSVC,
		logger:     logger,

		jobInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "info"),
			"Brigade job information.",
			[]string{"id", "build_id", "name", "image"}, nil,
		),
		jobStatusDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "status"),
			"Brigade job status.",
			[]string{"id", "status"}, nil,
		),
		jobDurationDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, jobSubSystem, "duration_seconds"),
			"Brigade job duration in seconds.",
			[]string{"id"}, nil,
		),
	}
}

// Collect satisfies subcollector.
func (j *job) Collect(ch chan<- prometheus.Metric) error {
	jobs, err := j.brigadeSVC.GetJobs()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		// Info metric.
		ch <- prometheus.MustNewConstMetric(
			j.jobInfoDesc,
			prometheus.GaugeValue,
			1,
			job.ID, job.BuildID, job.Name, job.Image)

		// Status metric.
		ch <- prometheus.MustNewConstMetric(
			j.jobStatusDesc,
			prometheus.GaugeValue,
			1,
			job.ID, job.Status)

		// Duration metric.
		// TODO: Think if it's 0 we should send the metric or not.
		ch <- prometheus.MustNewConstMetric(
			j.jobDurationDesc,
			prometheus.GaugeValue,
			job.Duration.Seconds(),
			job.ID)
	}

	return nil
}
