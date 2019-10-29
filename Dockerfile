FROM golang:1.12 AS builder

WORKDIR /go/src/github.com/controlplaneio/kubesec

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . .

RUN dep ensure -v -vendor-only \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec ./cmd/kubesec/*

# ===

FROM alpine:3.8

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/controlplaneio/kubesec/kubesec .
RUN chown -R app:app ./

COPY --from=stefanprodan/kubernetes-json-schema:latest /schemas/master-standalone /schemas/kubernetes-json-schema/master/master-standalone
RUN chown -R app:app /schemas

USER app

ENTRYPOINT ["./kubesec"]
CMD ["http", "8080"]
