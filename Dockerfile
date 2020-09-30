FROM golang:1.13 AS builder

WORKDIR /kubesec

COPY . .

WORKDIR /kubesec/v2

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec ./cmd/kubesec

# ===

FROM alpine:3.11

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=builder /kubesec/v2/kubesec .
COPY --from=stefanprodan/kubernetes-json-schema:latest /schemas/master-standalone /schemas/master-standalone-strict

RUN chown -R app:app ./ /schemas

USER app

ENTRYPOINT ["./kubesec"]
CMD ["http", "8080"]
