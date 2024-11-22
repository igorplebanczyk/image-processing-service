FROM postgres:latest

ARG PORT
ARG USER
ARG PASSWORD
ARG DB

ENV PORT=${PORT}
ENV POSTGRES_USER=${USER}
ENV POSTGRES_PASSWORD=${PASSWORD}
ENV POSTGRES_DB=${DB}

COPY ./sql/schema /docker-entrypoint-initdb.d/
COPY ./sql/objects /docker-entrypoint-initdb.d/

HEALTHCHECK --interval=10s --retries=5 --start-period=30s --timeout=5s \
  CMD pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}

EXPOSE ${PORT}