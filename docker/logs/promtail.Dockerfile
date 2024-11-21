FROM grafana/promtail:latest

COPY ./docker/logs/promtail.yaml /etc/promtail/promtail.yaml

EXPOSE 9080

CMD ["-config.file=/etc/promtail/promtail.yaml"]