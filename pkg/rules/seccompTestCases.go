package rules

// Currently, seccompProfile can appear under the following places (search for 'seccomp' in
// the PSS docs https://kubernetes.io/docs/concepts/security/pod-security-standards/):
//
//	spec.securityContext.seccompProfile.type
//	spec.containers[*].securityContext.seccompProfile.type
//	spec.initContainers[*].securityContext.seccompProfile.type
//	spec.ephemeralContainers[*].securityContext.seccompProfile.type

const (
	tcmanSeccompProfileMissing = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanSeccompProfileInSecCtxRD = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    seccompProfile:
      type: RuntimeDefault
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanSeccompProfileInSecCtxUn = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    seccompProfile:
      type: Unconfined
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanSeccompProfileInSecCtxLH = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    seccompProfile:
      type: LocalHost
      localhostProfile: profiles/audit.json
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanSeccompProfileInContainerRD = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: RuntimeDefault
`
	tcmanSeccompProfileInContainerUn = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: Unconfined
`
	tcmanSeccompProfileInContainerLH = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
	tcmanSeccompProfileInInitRD = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: RuntimeDefault
`
	tcmanSeccompProfileInInitUn = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: Unconfined
`
	tcmanSeccompProfileInInitLH = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
	tcmanSeccompProfileInEphRD = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: RuntimeDefault
`
	tcmanSeccompProfileInEphUn = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: Unconfined
`
	tcmanSeccompProfileInEphLH = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      seccompProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
)

type tcSeccompProfileType string

const (
	tcprofSeccompProfileMissing tcSeccompProfileType = ""
	tcprofSeccompUnconfined     tcSeccompProfileType = "Unconfined"
	tcprofSeccompRuntimeDefault tcSeccompProfileType = "RuntimeDefault"
	tcprofSeccompLocalhost      tcSeccompProfileType = "LocalHost"
)

var testCasesSeccomp = []struct {
	description         string
	expectedProfileType tcSeccompProfileType
	manifest            string
}{
	{
		description:         "Missing seccompProfile",
		expectedProfileType: tcprofSeccompProfileMissing,
		manifest:            tcmanSeccompProfileMissing,
	},
	{
		description:         "RuntimeDefault seccompProfile found inside .spec.securityContext",
		expectedProfileType: tcprofSeccompRuntimeDefault,
		manifest:            tcmanSeccompProfileInSecCtxRD,
	},
	{
		description:         "Unconfined seccompProfile found inside .spec.securityContext",
		expectedProfileType: tcprofSeccompUnconfined,
		manifest:            tcmanSeccompProfileInSecCtxUn,
	},
	{
		description:         "LocalHost seccompProfile found inside .spec.securityContext",
		expectedProfileType: tcprofSeccompLocalhost,
		manifest:            tcmanSeccompProfileInSecCtxLH,
	},
	{
		description:         "RuntimeDefault seccompProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofSeccompRuntimeDefault,
		manifest:            tcmanSeccompProfileInContainerRD,
	},
	{
		description:         "Unconfined seccompProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofSeccompUnconfined,
		manifest:            tcmanSeccompProfileInContainerUn,
	},
	{
		description:         "LocalHost seccompProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofSeccompLocalhost,
		manifest:            tcmanSeccompProfileInContainerLH,
	},
	{
		description:         "seccompProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofSeccompRuntimeDefault,
		manifest:            tcmanSeccompProfileInInitRD,
	},
	{
		description:         "Unconfined seccompProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofSeccompUnconfined,
		manifest:            tcmanSeccompProfileInInitUn,
	},
	{
		description:         "LocalHost seccompProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofSeccompLocalhost,
		manifest:            tcmanSeccompProfileInInitLH,
	},
	{
		description:         "RuntimeDefault seccompProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofSeccompRuntimeDefault,
		manifest:            tcmanSeccompProfileInEphRD,
	},
	{
		description:         "Unconfined seccompProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofSeccompUnconfined,
		manifest:            tcmanSeccompProfileInEphUn,
	},
	{
		description:         "LocalHost seccompProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofSeccompLocalhost,
		manifest:            tcmanSeccompProfileInEphLH,
	},
}
