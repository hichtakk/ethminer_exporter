FROM        quay.io/prometheus/busybox:glibc
MAINTAINER  Hirotaka Ichikawa <hichtakk@gmail.com>

COPY bin/ethminer_exporter /bin/ethminer_exporter

EXPOSE      8555
USER        nobody
ENTRYPOINT  [ "/bin/ethminer_exporter" ]