FROM redis:latest

ARG PORT
ENV PORT=${PORT}

COPY ./config/redis/redis.conf /usr/local/etc/redis/redis.conf

HEALTHCHECK --interval=10s --retries=5 --start-period=30s --timeout=5s \
  CMD ["redis-cli", "ping"]

EXPOSE ${PORT}