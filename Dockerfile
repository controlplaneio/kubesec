FROM golang:1.11 AS builder

RUN mkdir -p /go/src/github.com/sublimino/kubesec/

WORKDIR /go/src/github.com/sublimino/kubesec

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec ./cmd/kubesec/*

FROM alpine:3.8

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/sublimino/kubesec/kubesec .

RUN chown -R app:app ./

USER app

ENTRYPOINT ["./kubesec"]
CMD ["http", "9090"]
