FROM golang:1.16 as builder
COPY . /go/src/github.com/newrelic/nri-redis/
RUN cd /go/src/github.com/newrelic/nri-redis && \
    CGO_ENABLED=0 make compile && \
    strip ./bin/nri-redis

FROM alpine:latest

COPY --from=builder /go/src/github.com/newrelic/nri-redis/bin/nri-redis /nri-redis
COPY --from=builder /go/src/github.com/newrelic/nri-redis/src/spec.yaml /src/spec.yaml

EXPOSE 8080

CMD ["/nri-redis"]

