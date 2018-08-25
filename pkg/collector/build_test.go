package collector_test

import (
	"context"
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
	buildInfoDesc     = `Desc{fqName: "brigade_build_info", help: "Brigade build information.", constLabels: {}, variableLabels: [id project_id event_type provider version]}`
	buildStatusDesc   = `Desc{fqName: "brigade_build_status", help: "Brigade build status.", constLabels: {}, variableLabels: [id status]}`
	buildDurationDesc = `Desc{fqName: "brigade_build_duration_seconds", help: "Brigade build duration in seconds.", constLabels: {}, variableLabels: [id]}`
)

func TestBuildSubcollector(t *testing.T) {
	tests := []struct {
		name       string
		builds     []*brigade.Build
		expMetrics []metricResult
	}{
		{
			name: "With multiple builds the collected metrics should be of all the builds.",
			builds: []*brigade.Build{
				&brigade.Build{ID: "id1", ProjectID: "prj1", Type: "push", Provider: "gitlab", Version: "1234567890", Status: "Running", Duration: 125 * time.Second},
				&brigade.Build{ID: "id2", ProjectID: "prj2", Type: "pull_request", Provider: "github", Version: "1234567891", Status: "Pending", Duration: 340 * time.Second},
				&brigade.Build{ID: "id3", ProjectID: "prj3", Type: "deploy", Provider: "toilet", Version: "1234567892", Status: "Failed", Duration: 18 * time.Second},
			},
			expMetrics: []metricResult{
				metricResult{
					desc:       buildInfoDesc,
					labels:     labelMap{"id": "id1", "project_id": "prj1", "event_type": "push", "provider": "gitlab", "version": "1234567890"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildStatusDesc,
					labels:     labelMap{"id": "id1", "status": "Running"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildDurationDesc,
					labels:     labelMap{"id": "id1"},
					value:      125,
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       buildInfoDesc,
					labels:     labelMap{"id": "id2", "project_id": "prj2", "event_type": "pull_request", "provider": "github", "version": "1234567891"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildStatusDesc,
					labels:     labelMap{"id": "id2", "status": "Pending"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildDurationDesc,
					labels:     labelMap{"id": "id2"},
					value:      340,
					metricType: dto.MetricType_GAUGE,
				},

				metricResult{
					desc:       buildInfoDesc,
					labels:     labelMap{"id": "id3", "project_id": "prj3", "event_type": "deploy", "provider": "toilet", "version": "1234567892"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildStatusDesc,
					labels:     labelMap{"id": "id3", "status": "Failed"},
					value:      1,
					metricType: dto.MetricType_GAUGE,
				},
				metricResult{
					desc:       buildDurationDesc,
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
			mbsvc.On("GetBuilds").Once().Return(test.builds, nil)

			clr := collector.NewBuild(mbsvc, log.Dummy)

			ch := make(chan prometheus.Metric)

			go func() {
				clr.Collect(context.TODO(), ch)
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
