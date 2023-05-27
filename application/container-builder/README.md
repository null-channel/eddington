





# how to build a container image
Make sure the user running dagger has appropriate permissions (`rw`) on the socket file.

To use a buildkit daemon listening on TCP port `1234` on localhost:

```shell
export BUILDKIT_HOST=tcp://localhost:1234
```