FROM ubuntu

LABEL maintainer="dev@flipt.io"
LABEL org.opencontainers.image.name="flipt"
LABEL org.opencontainers.image.source="https://github.com/flipt-io/flipt"

RUN apt update && apt install -y postgresql-client \
    openssl \
    ca-certificates 

RUN mkdir -p /etc/flipt && \
    mkdir -p /var/opt/flipt

COPY flipt /
COPY config/migrations/ /etc/flipt/config/migrations/
COPY config/*.yml /etc/flipt/config/


EXPOSE 8080
EXPOSE 9000


CMD ["./flipt"]