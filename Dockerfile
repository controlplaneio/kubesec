FROM golang:1.19 AS downloader

ARG K8S_SCHEMA_VER=master

WORKDIR /schemas

RUN set -x && \
    if [ "${K8S_SCHEMA_VER}" != "master" ]; then K8S_SCHEMA_VER="v${K8S_SCHEMA_VER}"; fi && \
    BASE_URL="https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master" && \
    SCHEMA_PATH="${K8S_SCHEMA_VER}-standalone-strict" && \
    mkdir "${SCHEMA_PATH}" && \
    cd "${SCHEMA_PATH}" && \
    curl -sSL -O "${BASE_URL}/${SCHEMA_PATH}/pod-v1.json" && \
    curl -sSL -O "${BASE_URL}/${SCHEMA_PATH}/daemonset-apps-v1.json" && \
    curl -sSL -O "${BASE_URL}/${SCHEMA_PATH}/deployment-apps-v1.json" && \
    curl -sSL -O "${BASE_URL}/${SCHEMA_PATH}/statefulset-apps-v1.json"

FROM golang:1.19 AS builder

WORKDIR /kubesec

COPY main.go go.mod go.sum ./
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec .

# ===

FROM alpine:3.16

ARG K8S_SCHEMA_VER
ENV K8S_SCHEMA_VER=${K8S_SCHEMA_VER:-}
ENV SCHEMA_LOCATION=/schemas

RUN addgroup -S kubesec \
  && adduser -S -g kubesec kubesec \
  && apk --no-cache add ca-certificates

WORKDIR /home/kubesec

# This directory must follow the same structure ($SCHEMA_PATH) as the upstream
# schema location: github.com/yannh/kubernetes-json-schema
COPY --from=downloader /schemas /schemas
COPY --from=builder /kubesec/kubesec /bin/kubesec
COPY ./templates/ /templates

RUN chown -R kubesec:kubesec ./ /schemas /templates

USER kubesec

ENTRYPOINT ["kubesec"]
CMD ["http", "8080"]
