FROM golang:1.15 AS builder

WORKDIR /kubesec

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec .

# ===

FROM alpine:3.13.1

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY ./templates/ templates
COPY --from=builder /kubesec/kubesec .
COPY --from=stefanprodan/kubernetes-json-schema:latest /schemas/master-standalone /schemas/master-standalone-strict

RUN chown -R app:app ./ /schemas

USER app

ENTRYPOINT ["./kubesec"]
CMD ["http", "8080"]
