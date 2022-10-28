# Changelog

All notable changes to this project will be documented in this file.

## Table of Contents

- [2.12.0](#2120)
- [2.11.4](#2114)
- [2.11.3](#2113)
- [2.11.2](#2112)
- [2.11.1](#2111)
- [2.11.0](#2110)
- [2.10.2](#2102)
- [2.10.1](#2101)
- [2.10.0](#2100)
- [2.9.0](#290)
- [2.8.0](#280)
- [2.7.2](#272)
- [2.7.1](#271)
- [2.7.0](#270)
- [2.6.0](#260)
- [2.5.0](#250)
- [2.4.0](#240)
- [2.3.1](#231)
- [2.3.0](#230)
- [2.2.0](#220)
- [2.1.0](#210)
- [2.0.0](#200)
- [1.0.0](#100)

---

## `2.12.0`

- Update kubesec dependencies
- Update actions
- Migrate from kubeval to kubeconform
- Fix StatefulSet and VolumeClaimTemplate issues

## `2.11.4`

- Fix container builds so all tags are correctly built
- Split release and container release so they can be re-ran separately
- Bump dependencies

## `2.11.3`

- Bump dependencies
- Minor doc cleanup

## `2.11.2`

- Allow specifying schema location with `--schema-dir`
  - thanks @AndreasMili
- Fix LimitsMemory rule incorrectly using the RequestsLimit rule
  - thanks @AndreasMili

## `2.11.1`

- Split out actions so they can run only when necessary
- Bump dependencies
  - Includes a couple more breaking updates that required some additional work to integrate

## `2.11.0`

- Move assets in the containers to make them easier to access
- Fix changelog links
- Add exit-code override

## `2.10.2`

- drop ghcr until auth is fixed

## `2.10.1`

- actually push the container releases

## `2.10.0`

- add more release targets
- sunset i386 target
- add template directory to the Dockerfiles
- build and push containers on release
  - Docker Hub
  - GitHub Container Registry

## `2.9.0`

- add templating output format
- add provided sarif template
- add output location
- make go install and build easier by splitting cmd and a main.go in the root
- cleaned up docs
- made tests less brittle
- fix scratch container

## `2.8.0`

- fix issues processing multi doc yaml with empty elements
- added some more kubesec scan examples
- added the file name to the kubeval input
- added a flag to show the absolute filename instead

## `2.7.2`

- bump go and alpine versions
  - this is also part of making `go mod` happy with `v2`

## `2.7.1`

- further fixes to make `go mod` happy with `v2`
  - should resolve issues with tools that use `go list ./...` at the project root

## `2.7.0`

- fix go mod issues with `v2`
  - can use `go get` again

## `2.6.0`

- allow for piping into `kubesec scan` using `-` or `/dev/stdin`
  - `cat somefile.yml | kubesec scan -`
  - `cat somefile.yml | kubesec scan /dev/stdin`

## `2.5.0`

- improved in-toto integration

## `2.4.0`

- added passed to the JSON output
- note: repo tests now require `jq` - **only concerns maintainers**

## `2.3.1`

- patch to accept form data from the <https://kubesec.io> webpage sample form

## `2.3.0`

- moved everything to go modules

## `2.2.0`

- added in-toto support

## `2.1.0`

- add rule for `allowPrivilegeEscalation: true` with a score of -7
- add `points` field to each recommendation so the values that comprise the total score can be seen
- fix case sensitivity bug in `.capabilities.drop | index("ALL")`
- rules in `critical` and `advise` lists prioritised and returned in same order across runs

## `2.0.0`

- first open source release
- passes same acceptance tests as Kubesec v1
- more stringent analysis: scoring for a rule is multiplied by number of matches (previously the score was only applied
  once), initContainers are included in score, new securityContext directive support, seccomp and apparmor pod-targeting
  tighter
- CLI and HTTP server bundled in single binary

## `1.0.0`

- initial release at <https://kubesec.io>
- closed source
