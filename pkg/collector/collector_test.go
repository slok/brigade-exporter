package collector_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	mbrigade "github.com/slok/brigade-exporter/mocks/service/brigade"
	"github.com/slok/brigade-exporter/pkg/collector"
	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

var (
	t1 = time.Now()
	t2 = t1.Add(265 * time.Second)
	t3 = t2.Add(12 * time.Minute)
	t4 = t3.Add(1 * time.Hour)

	testProjects = []*brigade.Project{
		&brigade.Project{ID: "id1", Name: "Name1", Repository: "repo1", Namespace: "ns1", Worker: "worker1"},
		&brigade.Project{ID: "id2", Name: "Name2", Repository: "repo2", Namespace: "ns2", Worker: "worker2"},
		&brigade.Project{ID: "id3", Name: "Name3", Repository: "repo3", Namespace: "ns3", Worker: "worker3"},
	}
	testBuilds = []*brigade.Build{
		&brigade.Build{ID: "id1", ProjectID: "prj1", Type: "push", Provider: "gitlab", Version: "1234567890", Status: "Running", Duration: 125 * time.Second},
		&brigade.Build{ID: "id2", ProjectID: "prj2", Type: "pull_request", Provider: "github", Version: "1234567891", Status: "Pending", Duration: 340 * time.Second},
		&brigade.Build{ID: "id3", ProjectID: "prj3", Type: "deploy", Provider: "toilet", Version: "1234567892", Status: "Failed", Duration: 18 * time.Second},
	}
	testJobs = []*brigade.Job{
		&brigade.Job{ID: "id1", BuildID: "bld1", Name: "id-name-1", Image: "image1", Status: "Running", Duration: 125 * time.Second, Creation: t1, Start: t2},
		&brigade.Job{ID: "id2", BuildID: "bld2", Name: "id-name-2", Image: "image2", Status: "Pending", Duration: 340 * time.Second, Creation: t2, Start: t3},
		&brigade.Job{ID: "id3", BuildID: "bld3", Name: "id-name-3", Image: "image3", Status: "Failed", Duration: 18 * time.Second},
	}
)

func TestExporter(t *testing.T) {

	tests := []struct {
		name          string
		exporterCfg   collector.Config
		projects      []*brigade.Project
		builds        []*brigade.Build
		jobs          []*brigade.Job
		expMetrics    []string
		notExpMetrics []string
	}{
		{
			name:     "By default the exporter should return metrics from all subcollectors.",
			projects: testProjects,
			builds:   testBuilds,
			jobs:     testJobs,
			expMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="projects"} 1`,
				`brigade_exporter_collector_success{collector="builds"} 1`,
				`brigade_exporter_collector_success{collector="jobs"} 1`,

				// Brigade projects metrics.
				`brigade_project_info{id="id1",name="Name1",namespace="ns1",repository="repo1",worker="worker1"} 1`,
				`brigade_project_info{id="id2",name="Name2",namespace="ns2",repository="repo2",worker="worker2"} 1`,
				`brigade_project_info{id="id3",name="Name3",namespace="ns3",repository="repo3",worker="worker3"} 1`,

				// Brigade builds metrics.
				`brigade_build_info{event_type="deploy",id="id3",project_id="prj3",provider="toilet",version="1234567892"} 1`,
				`brigade_build_info{event_type="pull_request",id="id2",project_id="prj2",provider="github",version="1234567891"} 1`,
				`brigade_build_info{event_type="push",id="id1",project_id="prj1",provider="gitlab",version="1234567890"} 1`,
				`brigade_build_duration_seconds{id="id1"} 125`,
				`brigade_build_duration_seconds{id="id2"} 340`,
				`brigade_build_duration_seconds{id="id3"} 18`,
				`brigade_build_status{id="id1",status="Running"} 1`,
				`brigade_build_status{id="id2",status="Pending"} 1`,
				`brigade_build_status{id="id3",status="Failed"} 1`,

				// Brigade Jobs metrics.
				`brigade_job_info{build_id="bld1",id="id1",image="image1",name="id-name-1"} 1`,
				`brigade_job_info{build_id="bld2",id="id2",image="image2",name="id-name-2"} 1`,
				`brigade_job_info{build_id="bld3",id="id3",image="image3",name="id-name-3"} 1`,
				`brigade_job_duration_seconds{id="id1"} 125`,
				`brigade_job_duration_seconds{id="id2"} 340`,
				`brigade_job_duration_seconds{id="id3"} 18`,
				`brigade_job_status{id="id1",status="Running"} 1`,
				`brigade_job_status{id="id2",status="Pending"} 1`,
				`brigade_job_status{id="id3",status="Failed"} 1`,
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id1"}`, t1),
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id2"}`, t2),
				`brigade_job_create_time_seconds{id="id3"} 0`,
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id1"}`, t2),
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id2"}`, t3),
				`brigade_job_start_time_seconds{id="id3"} 0`,
			},
			notExpMetrics: []string{},
		},
		{
			name: "Disabling all the subcollectors except the projects it shouñld return only the projects metrics.",
			exporterCfg: collector.Config{
				DisableBuilds: true,
				DisableJobs:   true,
			},
			projects: testProjects,
			builds:   testBuilds,
			jobs:     testJobs,
			expMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="projects"} 1`,

				// Brigade projects metrics.
				`brigade_project_info{id="id1",name="Name1",namespace="ns1",repository="repo1",worker="worker1"} 1`,
				`brigade_project_info{id="id2",name="Name2",namespace="ns2",repository="repo2",worker="worker2"} 1`,
				`brigade_project_info{id="id3",name="Name3",namespace="ns3",repository="repo3",worker="worker3"} 1`,
			},
			notExpMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="builds"} 1`,
				`brigade_exporter_collector_success{collector="jobs"} 1`,

				// Brigade builds metrics.
				`brigade_build_info{event_type="deploy",id="id3",project_id="prj3",provider="toilet",version="1234567892"} 1`,
				`brigade_build_info{event_type="pull_request",id="id2",project_id="prj2",provider="github",version="1234567891"} 1`,
				`brigade_build_info{event_type="push",id="id1",project_id="prj1",provider="gitlab",version="1234567890"} 1`,
				`brigade_build_duration_seconds{id="id1"} 125`,
				`brigade_build_duration_seconds{id="id2"} 340`,
				`brigade_build_duration_seconds{id="id3"} 18`,
				`brigade_build_status{id="id1",status="Running"} 1`,
				`brigade_build_status{id="id2",status="Pending"} 1`,
				`brigade_build_status{id="id3",status="Failed"} 1`,

				// Brigade Jobs metrics.
				`brigade_job_info{build_id="bld1",id="id1",image="image1",name="id-name-1"} 1`,
				`brigade_job_info{build_id="bld2",id="id2",image="image2",name="id-name-2"} 1`,
				`brigade_job_info{build_id="bld3",id="id3",image="image3",name="id-name-3"} 1`,
				`brigade_job_duration_seconds{id="id1"} 125`,
				`brigade_job_duration_seconds{id="id2"} 340`,
				`brigade_job_duration_seconds{id="id3"} 18`,
				`brigade_job_status{id="id1",status="Running"} 1`,
				`brigade_job_status{id="id2",status="Pending"} 1`,
				`brigade_job_status{id="id3",status="Failed"} 1`,
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id1"}`, t1),
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id2"}`, t2),
				`brigade_job_create_time_seconds{id="id3"} 0`,
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id1"}`, t2),
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id2"}`, t3),
				`brigade_job_start_time_seconds{id="id3"} 0`,
			},
		},
		{
			name: "Disabling project subcollectors should return everything except project metrics.",
			exporterCfg: collector.Config{
				DisableProjects: true,
			},
			projects: testProjects,
			builds:   testBuilds,
			jobs:     testJobs,
			expMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="builds"} 1`,
				`brigade_exporter_collector_success{collector="jobs"} 1`,

				// Brigade builds metrics.
				`brigade_build_info{event_type="deploy",id="id3",project_id="prj3",provider="toilet",version="1234567892"} 1`,
				`brigade_build_info{event_type="pull_request",id="id2",project_id="prj2",provider="github",version="1234567891"} 1`,
				`brigade_build_info{event_type="push",id="id1",project_id="prj1",provider="gitlab",version="1234567890"} 1`,
				`brigade_build_duration_seconds{id="id1"} 125`,
				`brigade_build_duration_seconds{id="id2"} 340`,
				`brigade_build_duration_seconds{id="id3"} 18`,
				`brigade_build_status{id="id1",status="Running"} 1`,
				`brigade_build_status{id="id2",status="Pending"} 1`,
				`brigade_build_status{id="id3",status="Failed"} 1`,

				// Brigade Jobs metrics.
				`brigade_job_info{build_id="bld1",id="id1",image="image1",name="id-name-1"} 1`,
				`brigade_job_info{build_id="bld2",id="id2",image="image2",name="id-name-2"} 1`,
				`brigade_job_info{build_id="bld3",id="id3",image="image3",name="id-name-3"} 1`,
				`brigade_job_duration_seconds{id="id1"} 125`,
				`brigade_job_duration_seconds{id="id2"} 340`,
				`brigade_job_duration_seconds{id="id3"} 18`,
				`brigade_job_status{id="id1",status="Running"} 1`,
				`brigade_job_status{id="id2",status="Pending"} 1`,
				`brigade_job_status{id="id3",status="Failed"} 1`,
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id1"}`, t1),
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id2"}`, t2),
				`brigade_job_create_time_seconds{id="id3"} 0`,
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id1"}`, t2),
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id2"}`, t3),
				`brigade_job_start_time_seconds{id="id3"} 0`,
			},
			notExpMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="projects"} 1`,

				// Brigade projects metrics.
				`brigade_project_info{id="id1",name="Name1",namespace="ns1",repository="repo1",worker="worker1"} 1`,
				`brigade_project_info{id="id2",name="Name2",namespace="ns2",repository="repo2",worker="worker2"} 1`,
				`brigade_project_info{id="id3",name="Name3",namespace="ns3",repository="repo3",worker="worker3"} 1`,
			},
		},
		{
			name:     "A timeout in subcollectors should return bad collector success.",
			projects: testProjects,
			builds:   testBuilds,
			jobs:     testJobs,
			exporterCfg: collector.Config{
				CollectTimeout: 1, // 1 nanosecond is almost a timeout.
			},
			expMetrics: []string{
				// Exporter metrics.
				`brigade_exporter_collector_success{collector="projects"} 0`,
				`brigade_exporter_collector_success{collector="builds"} 0`,
				`brigade_exporter_collector_success{collector="jobs"} 0`,
			},
			notExpMetrics: []string{
				// Brigade projects metrics.
				`brigade_project_info{id="id1",name="Name1",namespace="ns1",repository="repo1",worker="worker1"} 1`,
				`brigade_project_info{id="id2",name="Name2",namespace="ns2",repository="repo2",worker="worker2"} 1`,
				`brigade_project_info{id="id3",name="Name3",namespace="ns3",repository="repo3",worker="worker3"} 1`,

				// Brigade builds metrics.
				`brigade_build_info{event_type="deploy",id="id3",project_id="prj3",provider="toilet",version="1234567892"} 1`,
				`brigade_build_info{event_type="pull_request",id="id2",project_id="prj2",provider="github",version="1234567891"} 1`,
				`brigade_build_info{event_type="push",id="id1",project_id="prj1",provider="gitlab",version="1234567890"} 1`,
				`brigade_build_duration_seconds{id="id1"} 125`,
				`brigade_build_duration_seconds{id="id2"} 340`,
				`brigade_build_duration_seconds{id="id3"} 18`,
				`brigade_build_status{id="id1",status="Running"} 1`,
				`brigade_build_status{id="id2",status="Pending"} 1`,
				`brigade_build_status{id="id3",status="Failed"} 1`,

				// Brigade Jobs metrics.
				`brigade_job_info{build_id="bld1",id="id1",image="image1",name="id-name-1"} 1`,
				`brigade_job_info{build_id="bld2",id="id2",image="image2",name="id-name-2"} 1`,
				`brigade_job_info{build_id="bld3",id="id3",image="image3",name="id-name-3"} 1`,
				`brigade_job_duration_seconds{id="id1"} 125`,
				`brigade_job_duration_seconds{id="id2"} 340`,
				`brigade_job_duration_seconds{id="id3"} 18`,
				`brigade_job_status{id="id1",status="Running"} 1`,
				`brigade_job_status{id="id2",status="Pending"} 1`,
				`brigade_job_status{id="id3",status="Failed"} 1`,
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id1"}`, t1),
				getUnixTimeMetric(`brigade_job_create_time_seconds{id="id2"}`, t2),
				`brigade_job_create_time_seconds{id="id3"} 0`,
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id1"}`, t2),
				getUnixTimeMetric(`brigade_job_start_time_seconds{id="id2"}`, t3),
				`brigade_job_start_time_seconds{id="id3"} 0`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mbsvc := &mbrigade.Interface{}
			mbsvc.On("GetProjects").Once().Return(test.projects, nil)
			mbsvc.On("GetBuilds").Once().Return(test.builds, nil)
			mbsvc.On("GetJobs").Once().Return(test.jobs, nil)

			// Create the exporter.
			clr := collector.NewExporter(test.exporterCfg, mbsvc, log.Dummy)
			promReg := prometheus.NewRegistry()
			promReg.MustRegister(clr)
			h := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})

			// Make request to ask for metrics.
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/metrics", nil)
			h.ServeHTTP(rec, req)

			resp := rec.Result()

			if assert.Equal(http.StatusOK, resp.StatusCode) {
				body, _ := ioutil.ReadAll(resp.Body)

				// Check all exp metrics are present.
				for _, expMetric := range test.expMetrics {
					assert.Contains(string(body), expMetric, "metric not present on the result of metrics service")
				}

				// Check all not exp metrics are not present.
				for _, notExpMetric := range test.notExpMetrics {
					assert.NotContains(string(body), notExpMetric, "metric present on the result of metrics service, it shouldn't")
				}
			}
		})
	}
}

func getUnixTimeMetric(metric string, t time.Time) string {
	return fmt.Sprintf(`%s %g`, metric, float64(t.Unix()))
}
