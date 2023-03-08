FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/tokend/course-certificates/sbt-bot
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/sbt-bot /go/src/gitlab.com/tokend/course-certificates/sbt-bot


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/sbt-bot /usr/local/bin/sbt-bot
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["sbt-bot"]
