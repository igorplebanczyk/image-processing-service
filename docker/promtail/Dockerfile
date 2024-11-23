FROM grafana/promtail:latest

COPY ./config/promtail/promtail.yaml /etc/promtail/promtail.yaml

EXPOSE 9080

CMD ["-config.file=/etc/promtail/promtail.yaml"]