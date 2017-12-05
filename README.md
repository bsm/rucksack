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
- 

With `_ "github.com/bsm/rucksack/met/datadog"` imported:
- `MET_DATADOG=datadog-token` (required)

With `_ "github.com/bsm/rucksack/met/runtime"` imported:
- `MET_RUNTIME=mem,heap,gc` (used set is equivalent to `all`)

Optional ENV:
- `HOST=hostname` (auto-detected)
- `PORT=8080` (optional)

## Licence

    Copyright 2017 Black Square Media Ltd

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
