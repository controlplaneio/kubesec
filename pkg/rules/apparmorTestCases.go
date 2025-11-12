package rules

// Currently, appArmorProfile can appear under the following places (search for 'AppArmor' in
// the PSS docs https://kubernetes.io/docs/concepts/security/pod-security-standards/):
//
// 	spec.securityContext.appArmorProfile.type
// 	spec.containers[*].securityContext.appArmorProfile.type
// 	spec.initContainers[*].securityContext.appArmorProfile.type
// 	spec.ephemeralContainers[*].securityContext.appArmorProfile

const (
	tcmanAppArmorProfileMissing = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanAppArmorProfileInSecCtxRD = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    appArmorProfile:
      type: RuntimeDefault
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanAppArmorProfileInSecCtxUn = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    appArmorProfile:
      type: Unconfined
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanAppArmorProfileInSecCtxLH = `
---
apiVersion: v1
kind: Pod
spec:
  securityContext:
    appArmorProfile:
      type: LocalHost
      localhostProfile: profiles/audit.json
  containers:
    - name: trustworthy-container
      image: sotrustworthy:latest
`
	tcmanAppArmorProfileInContainerRD = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: RuntimeDefault
`
	tcmanAppArmorProfileInContainerUn = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: Unconfined
`
	tcmanAppArmorProfileInContainerLH = `
---
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: trustworthy-container
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
	tcmanAppArmorProfileInInitRD = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: RuntimeDefault
`
	tcmanAppArmorProfileInInitUn = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: Unconfined
`
	tcmanAppArmorProfileInInitLH = `
---
apiVersion: v1
kind: Pod
spec:
  initContainers:
  - name: trustworthy-initcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
	tcmanAppArmorProfileInEphRD = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: RuntimeDefault
`
	tcmanAppArmorProfileInEphUn = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: Unconfined
`
	tcmanAppArmorProfileInEphLH = `
---
apiVersion: v1
kind: Pod
spec:
  ephemeralContainers:
  - name: trustworthy-ephcontainer
    image: sotrustworthy:latest
    securityContext:
      appArmorProfile:
        type: LocalHost
        localhostProfile: profiles/audit.json
`
)

type tcAppArmorProfileType string

const (
	tcprofAppArmorProfileMissing tcAppArmorProfileType = ""
	tcprofAppArmorUnconfined     tcAppArmorProfileType = "Unconfined"
	tcprofAppArmorRuntimeDefault tcAppArmorProfileType = "RuntimeDefault"
	tcprofAppArmorLocalhost      tcAppArmorProfileType = "LocalHost"
)

var testCasesApparmor = []struct {
	description         string
	expectedProfileType tcAppArmorProfileType
	manifest            string
}{
	{
		description:         "Missing appArmorProfile",
		expectedProfileType: tcprofAppArmorProfileMissing,
		manifest:            tcmanAppArmorProfileMissing,
	},
	{
		description:         "RuntimeDefault appArmorProfile found inside .spec.securityContext",
		expectedProfileType: tcprofAppArmorRuntimeDefault,
		manifest:            tcmanAppArmorProfileInSecCtxRD,
	},
	{
		description:         "Unconfined appArmorProfile found inside .spec.securityContext",
		expectedProfileType: tcprofAppArmorUnconfined,
		manifest:            tcmanAppArmorProfileInSecCtxUn,
	},
	{
		description:         "LocalHost appArmorProfile found inside .spec.securityContext",
		expectedProfileType: tcprofAppArmorLocalhost,
		manifest:            tcmanAppArmorProfileInSecCtxLH,
	},
	{
		description:         "RuntimeDefault appArmorProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofAppArmorRuntimeDefault,
		manifest:            tcmanAppArmorProfileInContainerRD,
	},
	{
		description:         "Unconfined appArmorProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofAppArmorUnconfined,
		manifest:            tcmanAppArmorProfileInContainerUn,
	},
	{
		description:         "LocalHost appArmorProfile found inside .spec.containers[*].securityContext",
		expectedProfileType: tcprofAppArmorLocalhost,
		manifest:            tcmanAppArmorProfileInContainerLH,
	},
	{
		description:         "RuntimeDefault appArmorProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorRuntimeDefault,
		manifest:            tcmanAppArmorProfileInInitRD,
	},
	{
		description:         "Unconfined appArmorProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorUnconfined,
		manifest:            tcmanAppArmorProfileInInitUn,
	},
	{
		description:         "LocalHost appArmorProfile found inside .spec.initContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorLocalhost,
		manifest:            tcmanAppArmorProfileInInitLH,
	},
	{
		description:         "RuntimeDefault appArmorProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorRuntimeDefault,
		manifest:            tcmanAppArmorProfileInEphRD,
	},
	{
		description:         "Unconfined appArmorProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorUnconfined,
		manifest:            tcmanAppArmorProfileInEphUn,
	},
	{
		description:         "LocalHost appArmorProfile found inside .spec.ephemeralContainers[*].securityContext",
		expectedProfileType: tcprofAppArmorLocalhost,
		manifest:            tcmanAppArmorProfileInEphLH,
	},
}
