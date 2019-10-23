# Rucksack

[![Build Status](https://travis-ci.org/bsm/rucksack.svg)](https://travis-ci.org/bsm/rucksack)

A collection of logging/instrumenting tools for building Go apps

## Basics

Package aims to simplify logging and instrumenting, so:

- package exposes functions, that work with default logger/instrumenter (which is usually not exposed)
- package is configured with ENV vars, like `LOG_TAGS=app:program go run program.go`

Recommended way to use:

- import and use `"github.com/bsm/rucksack/log"` or `"github.com/bsm/rucksack/met"` in non-main packages
- import extension packages like `_ "github.com/bsm/rucksack/met/datadog"` only in main package

## Logging

[![GoDoc](https://godoc.org/github.com/bsm/rucksack/log?status.svg)](https://godoc.org/github.com/bsm/rucksack/log)

```go
import "github.com/bsm/rucksack/log"
```

ENV:

- `LOG_NAME=projectname` (aliased as `APP_NAME`)
- `LOG_TAGS=foo:bar,baz:qux` (aliased as `APP_TAGS`)
- `LOG_LEVEL=INFO`
- `LOG_STACK=true` (any non-empty value will enable stack logging; this is an [expensive option](https://godoc.org/go.uber.org/zap#Stack))

Recommended way to use:

```go
package main

import "github.com/bsm/rucksack/log"

func main() {
  defer log.Sync()
  defer log.ErrorOnPanic()

  // do stuff
}
```

## Metrics

[![GoDoc](https://godoc.org/github.com/bsm/rucksack/met?status.svg)](https://godoc.org/github.com/bsm/rucksack/met)

```go
import (
  "github.com/bsm/rucksack/met"
  _ "github.com/bsm/rucksack/met/datadog"
  _ "github.com/bsm/rucksack/met/runtime"
)
```

ENV:

- `MET_NAME=projectname` (required; aliased as `APP_NAME`)
- `MET_TAGS=foo:bar,baz:qux` (aliased as `APP_TAGS`)

With `_ "github.com/bsm/rucksack/met/datadog"` imported:

- `MET_DATADOG=datadog-token` (required)
- `MET_DATADOG_DISABLE_COMPRESSION=true` (optional, disables compression when sending data to DataDog API)

With `_ "github.com/bsm/rucksack/met/runtime"` imported:

- `MET_RUNTIME=mem,heap,gc` (used set is equivalent to `all`)

Optional ENV:

- `HOST=hostname` (auto-detected)
- `PORT=8080` (optional)
