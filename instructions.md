<!-- omit in toc -->
# Panic with Go versions < 1.15.0

<!-- omit in toc -->
## Table of contents

- [Overview](#overview)
- [Go versions tested](#go-versions-tested)
- [Instructions](#instructions)
- [Output](#output)
  - [Passing](#passing)
  - [Failing](#failing)
- [References](#references)

## Overview

I've tested and confirmed that the steps in this doc reproduce a panic with Go
versions earlier than 1.15.0.

## Go versions tested

Panics:

- 1.13.14
- 1.13.15
- 1.14.0
- 1.14.1
- 1.14.2
- 1.14.3
- 1.14.4
- 1.14.5
- 1.14.6
- 1.14.7

Does not panic:

- 1.15beta1
- 1.15rc1
- 1.15rc2

## Instructions

1. ```cd /tmp```
1. ```git clone https://github.com/atc0005/go-teams-notify```
1. ```git checkout go-panic-repro```
1. ```docker run --rm -it -v `pwd`:`pwd` -w `pwd` golang:1.15.0 go test -v ./...```
1. ```docker run --rm -it -v `pwd`:`pwd` -w `pwd` golang:1.14.7 go test -v ./...```

Alternatively, use a local Go 1.14 installation (any patch release) and
replace the last two steps with `go test -v ./...`.

## Output

### Passing

```ShellSession
=== RUN   TestTeamsClientSend
    send_test.go:100: OK: test 0; test.error is of type *fmt.wrapError, err is of type *fmt.wrapError
    send_test.go:100: OK: test 1; test.error is of type *errors.errorString, err is of type *errors.errorString
2020/08/27 09:18:25 RoundTripper returned a response & error; ignoring response
    send_test.go:100: OK: test 2; test.error is of type *url.Error, err is of type *url.Error
2020/08/27 09:18:25 RoundTripper returned a response & error; ignoring response
    send_test.go:100: OK: test 3; test.error is of type *url.Error, err is of type *url.Error
    send_test.go:100: OK: test 4; test.error is of type *errors.errorString, err is of type *errors.errorString
--- PASS: TestTeamsClientSend (0.00s)
PASS
ok      github.com/atc0005/go-teams-notify/v2   0.005s
```

### Failing

```ShellSession
=== RUN   TestTeamsClientSend
    send_test.go:100: OK: test 0; test.error is of type *fmt.wrapError, err is of type *fmt.wrapError
    send_test.go:100: OK: test 1; test.error is of type *errors.errorString, err is of type *errors.errorString
2020/08/27 09:18:39 RoundTripper returned a response & error; ignoring response
    send_test.go:100: OK: test 2; test.error is of type *url.Error, err is of type *url.Error
2020/08/27 09:18:39 RoundTripper returned a response & error; ignoring response
    send_test.go:100: OK: test 3; test.error is of type *url.Error, err is of type *url.Error
--- FAIL: TestTeamsClientSend (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
        panic: runtime error: invalid memory address or nil pointer dereference [recovered]
        panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x4d68bd]

goroutine 19 [running]:
testing.tRunner.func1.1(0x6e28e0, 0x99af10)
        /usr/local/go/src/testing/testing.go:988 +0x30d
testing.tRunner.func1(0xc0000ca480)
        /usr/local/go/src/testing/testing.go:991 +0x3f9
panic(0x6e28e0, 0x99af10)
        /usr/local/go/src/runtime/panic.go:969 +0x166
io/ioutil.readAll.func1(0xc00010b6c0)
        /usr/local/go/src/io/ioutil/ioutil.go:30 +0x101
panic(0x6e28e0, 0x99af10)
        /usr/local/go/src/runtime/panic.go:969 +0x166
bytes.(*Buffer).ReadFrom(0xc00010b648, 0x0, 0x0, 0xc0000932b0, 0x0, 0x0)
        /usr/local/go/src/bytes/buffer.go:204 +0x7d
io/ioutil.readAll(0x0, 0x0, 0x200, 0x0, 0x0, 0x0, 0x0, 0x0)
        /usr/local/go/src/io/ioutil/ioutil.go:36 +0xe3
io/ioutil.ReadAll(...)
        /usr/local/go/src/io/ioutil/ioutil.go:45
github.com/atc0005/go-teams-notify/v2.teamsClient.Send(0xc000097650, 0x74596e, 0x26, 0x73b510, 0xb, 0x741bdb, 0x1d, 0x0, 0x0, 0x0, ...)
        /tmp/go-teams-notify/send.go:70 +0x2ab
github.com/atc0005/go-teams-notify/v2.TestTeamsClientSend(0xc0000ca480)
        /tmp/go-teams-notify/send_test.go:89 +0x55e
testing.tRunner(0xc0000ca480, 0x7526e8)
        /usr/local/go/src/testing/testing.go:1039 +0xdc
created by testing.(*T).Run
        /usr/local/go/src/testing/testing.go:1090 +0x372
FAIL    github.com/atc0005/go-teams-notify/v2   0.014s
FAIL
```

## References

- Changes that triggered the panic behavior
  - <https://github.com/atc0005/go-teams-notify/pull/43/commits/6db6217e7daac6da6ed43dd7cf02032a22c84af2>

- Apparent "fix" for the panic
  - <https://github.com/atc0005/go-teams-notify/pull/43/commits/414b97a31efc543cbef48ed51847882e1fde55d7>
    - specifically [line 123](https://github.com/atc0005/go-teams-notify/pull/43/commits/414b97a31efc543cbef48ed51847882e1fde55d7#diff-f78e54225b3d916938acc1771ac7811aR123)

- Guide for testing a http client
  - <http://hassansin.github.io/Unit-Testing-http-client-in-Go>
  - credit: this is what helped me understand what was happening, what needed
    to be changed to resolve the issue
