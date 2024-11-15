FROM postgres:latest

ARG PORT
ARG USER
ARG DB

ENV PORT=${PORT}
ENV USER=${USER}
ENV DB=${DB}

COPY ./sql/schema /docker-entrypoint-initdb.d/

HEALTHCHECK --interval=10s --retries=5 --start-period=30s --timeout=5s \
  CMD pg_isready -U ${USER} -d ${DB}

EXPOSE ${PORT}