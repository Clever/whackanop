# whackanop

`whackanop` monitors mongodb for long-running operations, killing any that it finds.

## Motivation

See [this blog post](http://blog.mongolab.com/2014/02/mongodb-currentop-killop/).

## Usage

```bash
$ whackanop -h
Usage of whackanop:
  -debug=true: in debug mode, operations that match the query are logged instead of killed
  -interval=1: how often, in seconds, to poll mongo for operations
  -mongourl="localhost": mongo url to connect to
  -query="{\"op\": \"query\", \"secs_running\": {\"$gt\": 60}}": query sent to db.currentOp()
```

## Installation

`go get github.com/Clever/whackanop`

## Local Development

Set this repository up in the [standard location](https://golang.org/doc/code.html) in your `GOPATH`, i.e. `$GOPATH/src/github.com/Clever/whackanop`.
Once this is done, `make test` runs the tests.
