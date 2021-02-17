FROM golang:1.16.0 AS builder

WORKDIR /kubesec

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec .

# ===

FROM alpine:3.13.1

RUN addgroup -S kubesec \
    && adduser -S -g kubesec kubesec \
    && apk --no-cache add ca-certificates

WORKDIR /home/kubesec

COPY --from=builder /kubesec/kubesec /bin/kubesec
COPY --from=stefanprodan/kubernetes-json-schema:latest /schemas/master-standalone /schemas/master-standalone-strict
COPY ./templates/ /templates

RUN chown -R kubesec:kubesec ./ /schemas

USER kubesec

ENTRYPOINT ["kubesec"]
CMD ["http", "8080"]
