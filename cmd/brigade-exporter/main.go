package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	azurebrigade "github.com/Azure/brigade/pkg/storage/kube"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/slok/brigade-exporter/pkg/collector"
	"github.com/slok/brigade-exporter/pkg/log"
	"github.com/slok/brigade-exporter/pkg/service/brigade"
)

const (
	gracePeriod = 5 * time.Second
	versionFMT  = "brigade-exporter %s"
)

var (
	// Version is the app version.
	Version = "dev"
)

// Main is the main programm
type Main struct {
	flags  *flags
	logger log.Logger
}

// Run will run the main program.
func (m *Main) Run() error {

	if m.flags.version {
		m.printVersion()
		return nil
	}

	// If not development json logger.
	m.logger = log.Base(!m.flags.development)

	if m.flags.debug {
		m.logger.Set("debug")
	}

	var g run.Group

	// Signal capturing.
	{
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		g.Add(
			func() error {
				s := <-c
				m.logger.Infof("signal %s received", s)
				return nil
			},
			func(error) {},
		)
	}

	// Exporter.
	{
		// Prepare Services.
		var brigadeSVC brigade.Interface
		if m.flags.fake {
			brigadeSVC = brigade.NewFake()
			m.logger.Warnf("exporter running in faked mode")
		} else {
			k8scli, err := m.createKubernetesClient()
			if err != nil {
				return err
			}
			brigadeCli := azurebrigade.New(k8scli, m.flags.namespace)
			brigadeSVC = brigade.New(brigadeCli, m.logger)
		}

		// Prepare collector.
		clr := collector.NewExporter(collector.Config{}, brigadeSVC, m.logger)
		promReg := prometheus.NewRegistry()
		promReg.MustRegister(clr)

		// Prepare server
		h := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})
		mux := http.NewServeMux()
		mux.Handle(m.flags.metricsPath, h)
		s := http.Server{
			Handler: mux,
			Addr:    m.flags.listenAddress,
		}
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Brigade Exporter</title></head>
			<body>
			<h1>Brigade Exporter</h1>
			<p><a href="` + m.flags.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})

		g.Add(
			func() error {
				m.logger.Infof("listening on %s", m.flags.listenAddress)
				return s.ListenAndServe()
			},
			func(error) {
				ctx, _ := context.WithTimeout(context.Background(), gracePeriod)
				if err := s.Shutdown(ctx); err != nil {
					m.logger.Errorf("error while draining connections: %s", err)
				}
			},
		)

	}

	return g.Run()
}

// loadKubernetesConfig loads kubernetes configuration based on flags.
func (m *Main) loadKubernetesConfig() (*rest.Config, error) {
	var cfg *rest.Config
	// If devel mode then use configuration flag path.
	if m.flags.development {
		config, err := clientcmd.BuildConfigFromFlags("", m.flags.kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("could not load configuration: %s", err)
		}
		cfg = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading kubernetes configuration inside cluster, check app is running outside kubernetes cluster or run in development mode: %s", err)
		}
		cfg = config
	}

	return cfg, nil
}

func (m *Main) createKubernetesClient() (kubernetes.Interface, error) {
	config, err := m.loadKubernetesConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// printVersion prints the version of the app.
func (m *Main) printVersion() {
	fmt.Fprintf(os.Stdout, versionFMT, Version)
}

func main() {
	m := &Main{
		flags: NewFlags(),
	}

	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
