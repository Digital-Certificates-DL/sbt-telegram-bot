FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/tokend/course-certificates/sbt-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/sbt-svc /go/src/gitlab.com/tokend/course-certificates/sbt-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/sbt-svc /usr/local/bin/sbt-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["sbt-svc"]
