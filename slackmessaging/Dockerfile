FROM erikedin/nonaplumber as builder

COPY . /go/src/github.com/dandeliondeathray/nona/slackmessaging
RUN go get github.com/dandeliondeathray/nona/slackmessaging/cmd/slackmessaging
RUN go install github.com/dandeliondeathray/nona/slackmessaging/cmd/slackmessaging

FROM golang:1.8.4-stretch

COPY --from=builder /go/bin/slackmessaging /app/slackmessaging
