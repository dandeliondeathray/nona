ARG PLUMBER_TAG=latest
FROM erikedin/nonaplumber:${PLUMBER_TAG} as builder

COPY . /go/src/github.com/dandeliondeathray/nona/service/puzzlestore
RUN go get github.com/dandeliondeathray/nona/service/puzzlestore/cmd/nona_puzzlestore
RUN go install github.com/dandeliondeathray/nona/service/puzzlestore/cmd/nona_puzzlestore

FROM golang:1.8.4-stretch

COPY --from=builder /go/bin/nona_puzzlestore /app/nona_puzzlestore

CMD [ "/app/nona_puzzlestore" ]