# ethminer_exporter
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://hub.docker.com/r/hichtakk/ethminer_exporter/)

Prometheus exporter reporting ethminer hashrate  

# Setup

```sh
$ make build
```

# Usage

```sh
$ bin/ethminer_exporter
```

By default, ethminer_exporter starts listening on port 8555.
Then, visit http://localhost:8555/metrics?target=1.2.3.4:1234 where 1.2.3.4 is the IP of target ethmier host and 1234 is the JSON-RPC port ethminer listening. To enable JSON-RPC API for ethminer, you need to run with `--api-port` and specify listening port.
If they are specified correctly, ethminer_exporter reports its total hashrate at that time.

```
# HELP ethminer_totalhashrate Hashrate [H/s]
# TYPE ethminer_totalhashrate gauge
ethminer_totalhashrate 3.77168217e+08
```
