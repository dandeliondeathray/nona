FROM golang:1.8.3 as builder

RUN mkdir -p /go/src/github.com/dandeliondeathray/nona
COPY . /go/src/github.com/dandeliondeathray/nona
RUN go get github.com/dandeliondeathray/nona/...
RUN CGO_ENABLED=0 GOOS=linux go install -a -tags netgo -ldflags '-w' github.com/dandeliondeathray/nona/cmd/nona_slack

FROM scratch
LABEL maintainer="erikedin@users.noreply.github.com"
COPY --from=builder /go/bin/nona_slack /
