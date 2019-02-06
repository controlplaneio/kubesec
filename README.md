# kubesec

Validate the security parameters of Kubernetes YAML resources.

Currently supported types: Pod, Deployment, StatefulSet, DaemonSet

# WIP towards 1.0

1. This PR needs merging for goreleaser to work with signed commits https://github.com/goreleaser/goreleaser/pull/953
1. Jenkinsfile is a skeleton

```
make dep # install golang dependencies
make test-go # golang unit tests
make build # requires goreleaser, see PR above
```

