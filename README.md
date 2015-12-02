# whackanop

`whackanop` monitors mongodb for long-running operations, killing any that it finds.

## Motivation

See [this blog post](http://blog.mongolab.com/2014/02/mongodb-currentop-killop/).

Or this [StackExchange post](https://dba.stackexchange.com/questions/60029/how-do-i-safely-kill-long-running-operations-in-mongodb).

## Usage

```bash
$ whackanop -h
Usage of whackanop:
  -debug=true: in debug mode, operations that match the query are logged instead of killed
  -interval=1: how often, in seconds, to poll mongo for operations
  -mongourl="mongodb://localhost?connect=direct": mongo url to connect to. Must specify connect=direct to guarantee admin commands are run on the specified server.
  -query="{\"op\": \"query\", \"secs_running\": {\"$gt\": 60}}": query sent to db.currentOp()
  -verbose=false: more verbose logging
  -version=false: print the version and exit
```

## Installation

Install from source via `go get github.com/Clever/whackanop`, or download a release on the [releases](https://github.com/Clever/whackanop/releases) page.

## Local Development

Set this repository up in the [standard location](https://golang.org/doc/code.html) in your `GOPATH`, i.e. `$GOPATH/src/github.com/Clever/whackanop`.
Once this is done, `make test` runs the tests.

The release process requires a cross-compilation toolchain.
[`gox`](https://github.com/mitchellh/gox) can install the toolchain with one command: `gox -build-toolchain`.
From there you can build release tarballs for different OS and architecture combinations with `make release`.

### Rolling an official release

Official releases are listed on the [releases](https://github.com/Clever/whackanop/releases) page.
To create an official release:

1. On `master`, bump the version in the `VERSION` file in accordance with [semver](http://semver.org/).
You can do this with [`gitsem`](https://github.com/clever/gitsem), but make sure not to create the tag, e.g. `gitsem -tag=false patch`.

2. Push the change to Github. Drone will automatically create a release for you.

## Vendoring

Please view the [dev-handbook for instructions](https://github.com/Clever/dev-handbook/blob/master/golang/godep.md).
