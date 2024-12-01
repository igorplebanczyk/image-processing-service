FROM traefik:latest

CMD ["traefik", "--api.insecure=true", "--providers.docker=true", "--entrypoints.web.address=:80", "--accesslog=true"]

EXPOSE 80
