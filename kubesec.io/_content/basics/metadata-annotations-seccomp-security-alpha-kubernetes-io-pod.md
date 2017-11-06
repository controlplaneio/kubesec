Seccomp is a system call filtering facility in the Linux kernel which lets applications define limits on system calls they may make, and what should happen when system calls are made. Seccomp is used to reduce the attack surface available to applications.
<small>[source](https://github.com/kubernetes/kubernetes/blob/release-1.4/docs/design/seccomp.md)</small>

Specify a Seccomp profile for all containers of the Pod:

```
seccomp.security.alpha.kubernetes.io/pod
```

Specify a Seccomp profile for an individual container:

```
container.seccomp.security.alpha.kubernetes.io/${container_name}
```

## External Links
- [Seccomp Design doc](https://github.com/kubernetes/kubernetes/blob/release-1.4/docs/design/seccomp.md)
- [OCI Runtime Spec](https://github.com/opencontainers/runtime-spec/blob/master/config-linux.md#seccomp)
- [Seccomp filtering at Kernel.org](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)
- [Linux Seccomp examples](https://github.com/torvalds/linux/tree/master/samples/seccomp)
