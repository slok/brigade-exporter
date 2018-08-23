package collector_test

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	mbrigade "github.com/slok/brigade-exporter/mocks/service/brigade"
	"github.com/slok/brigade-exporter/pkg/collector"
	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	jobInfoDesc     = `Desc{fqName: "brigade_job_info", help: "Brigade job information.", constLabels: {}, variableLabels: [id build_id name image]}`
	jobStatusDesc   = `Desc{fqName: "brigade_job_status", help: "Brigade job status.", constLabels: {}, variableLabels: [id status]}`
	jobDurationDesc = `Desc{fqName: "brigade_job_duration_seconds", help: "Brigade job duration in seconds.", constLabels: {}, variableLabels: [id]}`
)

func TestJobSubcollector(t *testing.T) {
	tests := []struct {
		name        string
		exporterCfg collector.Config
		jobs        []*brigade.Job
		expMetrics  []metricResult
	}{
		{
			name: "With multiple jobs the collected metrics should be of all the jobs.",
			jobs: []*brigade.Job{
				&brigade.Job{ID: "id1", BuildID: "bld1", Name: "id-name-1", Image: "image1", Status: "Running", Duration: 125 * time.Second},
				&brigade.Job{ID: "id2", BuildID: "bld2", Name: "id-name-2", Image: "image2", Status: "Pending", Duration: 340 * time.Second},
				&brigade.Job{ID: "id3", BuildID: "bld3", Name: "id-name-3", Image: "image3", Status: "Failed", Duration: 18 * time.Second},
			},
			expMetrics: []metricResult{
				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id1", "build_id": "bld1", "name": "id-name-1", "image": "image1"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id1", "status": "Running"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id1"},
					value:      125,
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id2", "build_id": "bld2", "name": "id-name-2", "image": "image2"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id2", "status": "Pending"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id2"},
					value:      340,
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       jobInfoDesc,
					labels:     labelMap{"id": "id3", "build_id": "bld3", "name": "id-name-3", "image": "image3"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobStatusDesc,
					labels:     labelMap{"id": "id3", "status": "Failed"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       jobDurationDesc,
					labels:     labelMap{"id": "id3"},
					value:      18,
					metricType: dto.MetricType_GAUGE,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mbsvc := &mbrigade.Interface{}
			mbsvc.On("GetJobs").Once().Return(test.jobs, nil)

			clr := collector.NewJob(mbsvc, log.Dummy)

			ch := make(chan prometheus.Metric)

			go func() {
				clr.Collect(ch)
				close(ch)
			}()

			// Get the metrics
			var got []metricResult
			for m := range ch {
				got = append(got, readMetric(m))
			}

			// Check metrics are ok.
			assert.Equal(test.expMetrics, got)
		})
	}
}
