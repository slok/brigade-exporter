package main

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

// Defaults.
const (
	listenAddrDef  = ":9480"
	metricsPathDef = "/metrics"
	namespaceDef   = "default"
)

// flags are the flags of the app
type flags struct {
	fs *flag.FlagSet

	kubeConfig              string
	listenAddress           string
	metricsPath             string
	namespace               string
	disableProjectCollector bool
	disableBuildCollector   bool
	disableJobCollector     bool
	development             bool
	fake                    bool
	debug                   bool
	version                 bool
}

// newFlags returns a new flags object.
func newFlags() *flags {
	f := &flags{
		fs: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
	f.init()

	return f
}

func (f *flags) init() {
	kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// register flags
	f.fs.StringVar(&f.kubeConfig, "kubeconfig", kubehome, "kubernetes configuration path, only used when development mode enabled")
	f.fs.StringVar(&f.listenAddress, "listen-addr", listenAddrDef, "the address the exporter will be serving the metrics")
	f.fs.StringVar(&f.metricsPath, "metrics-path", metricsPathDef, "the path to serve the metrics")
	f.fs.StringVar(&f.namespace, "namespace", namespaceDef, "the namespace of brigade")
	f.fs.BoolVar(&f.disableProjectCollector, "disable-project-collector", false, "disables the metric gathering for brigade projects")
	f.fs.BoolVar(&f.disableBuildCollector, "disable-build-collector", false, "disables the metric gathering for brigade builds")
	f.fs.BoolVar(&f.disableJobCollector, "disable-job-collector", false, "disables the metric gathering for brigade jobs")
	f.fs.BoolVar(&f.development, "development", false, "development flag will run the exporter in development mode")
	f.fs.BoolVar(&f.fake, "fake", false, "fake flag will run the exporter faking the data from brigade")
	f.fs.BoolVar(&f.debug, "debug", false, "enable debug mode")
	f.fs.BoolVar(&f.version, "version", false, "show version")

	// Parse flags
	f.fs.Parse(os.Args[1:])
}
