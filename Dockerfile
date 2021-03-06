FROM ubuntu:21.04 AS builder

RUN apt update -y
RUN apt upgrade -y

RUN apt install -y locales
RUN apt install -y sudo

RUN echo "LC_ALL=en_US.UTF-8" >> /etc/environment && \
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen && \
    echo "LANG=en_US.UTF-8" > /etc/locale.conf && \
    locale-gen en_US.UTF-8

RUN useradd -m -G sudo developer
RUN echo 'developer:developer' | chpasswd
USER developer

RUN echo developer | sudo -S DEBIAN_FRONTEND="noninteractive" apt install -y golang
RUN echo developer | sudo -S apt install -y ca-certificates && sudo update-ca-certificates
RUN echo developer | sudo -S apt install -y make git vim protobuf-compiler

ENV GOPATH /home/developer/go
ENV PATH $PATH:/home/developer/go/bin


COPY . /home/developer/go/src/github.com/ozoncp/ocp-remind-api
RUN echo developer | sudo -S chown -R developer /home/developer/

WORKDIR /home/developer/go/src/github.com/ozoncp/ocp-remind-api

RUN make deps && make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
#ENV REMINDS_DB_URL postres://postgres:postgres@localhost:5432/reminds?sslmode=none
ENV REMINDS_CONF_FILENAME ../configuration.yaml
COPY --from=builder /home/developer/go/src/github.com/ozoncp/ocp-remind-api/bin/ocp-remind-api .
RUN chown root:root ocp-remind-api
EXPOSE 82
CMD ["./ocp-remind-api"]