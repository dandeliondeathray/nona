ARG PLUMBER_TAG=latest
FROM erikedin/nonaplumber:${PLUMBER_TAG} as builder

COPY . /go/src/github.com/dandeliondeathray/nona/service/puzzleresolver
RUN go get github.com/dandeliondeathray/nona/service/puzzleresolver/cmd/nona_puzzleresolver
RUN go install github.com/dandeliondeathray/nona/service/puzzleresolver/cmd/nona_puzzleresolver

FROM golang:1.8.4-stretch

COPY --from=builder /go/bin/nona_puzzleresolver /app/nona_puzzleresolver

CMD [ "/app/nona_puzzleresolver" ]