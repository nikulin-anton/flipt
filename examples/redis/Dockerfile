FROM markphelps/flipt:dev

USER root

RUN apt-get update && apt-get install -y git

RUN git clone https://github.com/vishnubob/wait-for-it.git /tmp && \
    chmod +x /tmp/wait-for-it.sh
