FROM golang:1.23 as builder

ENV LANG en_US.utf8
WORKDIR /logbench

COPY . .
RUN make

FROM ubuntu:22.04 as base

WORKDIR /logbench
COPY --from=builder /logbench/bin/logbench /usr/local/bin/

ENTRYPOINT ["logbench"]
