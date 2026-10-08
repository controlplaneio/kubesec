FROM index.docker.io/library/alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS downloader

ARG K8S_SCHEMA_VER=master

WORKDIR /schemas

RUN set -x && \
    apk add --no-cache curl && \
    if [ "${K8S_SCHEMA_VER}" != "master" ]; then K8S_SCHEMA_VER="v${K8S_SCHEMA_VER}"; fi && \
    BASE_URL="https://raw.githubusercontent.com/yannh/kubernetes-json-schema/master" && \
    SCHEMA_PATH="${K8S_SCHEMA_VER}-standalone-strict" && \
    mkdir "${SCHEMA_PATH}" && \
    curl -sSL --output-dir "${SCHEMA_PATH}" -O "${BASE_URL}/${SCHEMA_PATH}/pod-v1.json" && \
    curl -sSL --output-dir "${SCHEMA_PATH}" -O "${BASE_URL}/${SCHEMA_PATH}/daemonset-apps-v1.json" && \
    curl -sSL --output-dir "${SCHEMA_PATH}" -O "${BASE_URL}/${SCHEMA_PATH}/deployment-apps-v1.json" && \
    curl -sSL --output-dir "${SCHEMA_PATH}" -O "${BASE_URL}/${SCHEMA_PATH}/statefulset-apps-v1.json"

FROM index.docker.io/library/golang:1.27@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea AS builder

WORKDIR /kubesec

COPY main.go go.mod go.sum ./
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubesec .

# ===

FROM index.docker.io/library/alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

ARG K8S_SCHEMA_VER
ENV K8S_SCHEMA_VER=${K8S_SCHEMA_VER:-}
ENV SCHEMA_LOCATION=/schemas

RUN addgroup -S kubesec \
  && adduser -S -g kubesec kubesec \
  && apk --no-cache add ca-certificates

WORKDIR /home/kubesec

COPY --from=builder /kubesec/kubesec /bin/kubesec
COPY --chown=kubesec ./templates/ /templates
# This directory must follow the same structure ($SCHEMA_PATH) as the upstream
# schema location: github.com/yannh/kubernetes-json-schema
COPY --from=downloader --chown=kubesec /schemas /schemas

USER kubesec

ENTRYPOINT ["kubesec"]
CMD ["http", "8080"]
