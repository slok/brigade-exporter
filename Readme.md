# brigade-exporter

brigade-exporter is a Prometheus metrics exporter for [Brigade].

This exporter is designed to be run along with a brigade installation, if you have multiple brigades you will have multiple brigade-exporters, one per brigade installation. This follows the philosophy of prometheus exporters of one exporter per app instance.

## Run

There is already a docker image ready to run the exporter in `slok/brigade-exporter`. It has different options to run.

### Run outside the cluster

If you want to tests the exporter outside the cluster in a brigade installation you ca use `--development` flag. You will need kubectl configuration and the context to the correct cluster set.

```bash
docker run --rm -it -p 9477:9477 slok/brigade-exporter \
    --debug \
    --development \
    --namespace ${MY_BRIGADE_NAMESPACE}
```

go to http://127.0.0.1:9477/metrics

### Run in fake mode

If you are developing the exporter can fake a brigade installation and return fake data using `--fake` flag. This is very handy to develop.

## Deployment

TODO

### RBAC

TODO

## Metrics

### Exporter metrics

| Metric                                      | Type  | Meaning                            | Labels    |
| ------------------------------------------- | ----- | ---------------------------------- | --------- |
| brigade_exporter_collector_success          | gauge | Whether a collector succeeded      | collector |
| brigade_exporter_collector_duration_seconds | gauge | Collector time duration in seconds | collector |

### Project metrics

| Metric               | Type  | Meaning                     | Labels                                  |
| -------------------- | ----- | --------------------------- | --------------------------------------- |
| brigade_project_info | gauge | Brigade project information | id, name, namespace, repository, worker |

### Build metrics

| Metric                         | Type  | Meaning                           | Labels                                        |
| ------------------------------ | ----- | --------------------------------- | --------------------------------------------- |
| brigade_build_info             | gauge | Brigade build information         | id, project_id, event_type, provider, version |
| brigade_build_status           | gauge | Brigade build status              | id, status                                    |
| brigade_build_duration_seconds | gauge | Brigade build duration in seconds | id                                            |

### Job metrics

| Metric                       | Type  | Meaning                         | Labels                    |
| ---------------------------- | ----- | ------------------------------- | ------------------------- |
| brigade_job_info             | gauge | Brigade job information         | id, build_id, image, name |
| brigade_job_status           | gauge | Brigade job status              | id, status                |
| brigade_job_duration_seconds | gauge | Brigade job duration in seconds | id                        |

### Disabling metrics

You can disable metrics using flags.

- `--disable-project-collector`: Disables all the metrics of projects.
- `--disable-build-collector`: Disables all the metircs of builds.
- `--disable-job-collector`: Disables all the jobs metrics. If you have lots of jobs, this could improve the gathering and storage of metrics.

## Build from source

You can build your own brigade-exporter from source using:

```bash
make build-binary
```

to build the binary or

```bash
make build-image
```

to build the image.

[brigade]: https://brigade.sh
