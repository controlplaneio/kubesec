FROM golang:1.11 AS builder

RUN mkdir -p /go/src/github.com/sublimino/kubesec/

WORKDIR /go/src/github.com/sublimino/kubesec

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec ./cmd/kubesec/*

FROM golang:1.11 as schemas

RUN mkdir -p /schemas
WORKDIR /schemas
ADD https://github.com/garethr/kubernetes-json-schema/archive/master.tar.gz .
RUN tar xzf master.tar.gz --strip 1
RUN rm master.tar.gz

FROM alpine:3.8

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=builder /go/src/github.com/sublimino/kubesec/kubesec .
RUN chown -R app:app ./

COPY --from=schemas /schemas/master-standalone /schemas/kubernetes-json-schema/master/master-standalone
RUN chown -R app:app /schemas

USER app

ENTRYPOINT ["./kubesec"]
CMD ["http", "9090"]
