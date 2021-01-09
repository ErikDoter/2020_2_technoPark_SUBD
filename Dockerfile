FROM golang:1.13-stretch AS builder

WORKDIR /build
COPY . .

USER root
RUN go build  ./cmd/apiserver/run.go

FROM ubuntu:20.04
COPY . .

EXPOSE 5432
EXPOSE 5000

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update && apt -y install mysql

USER mysql

RUN  mysql -u root &&\
    create database forum &&\
    use forum &&\
     &&\
    exit;

#USER postgres
#CMD ["/usr/lib/postgresql/12/bin/postgres", "-D", "/var/lib/postgresql/12/main", "-c", "config_file=/etc/postgresql/12/main/postgresql.conf"]
#CMD ./run

USER root
COPY --from=builder  /build/run /usr/bin
CMD /etc/init.d/postgresql start && run