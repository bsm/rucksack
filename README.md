# Rucksack

[![Build Status](https://travis-ci.org/bsm/rucksack.svg)](https://travis-ci.org/bsm/rucksack)

A collection of logging/instrumenting tools for building Go apps

## Basics

Package aims to simplify logging and instrumenting, so:
- package exposes functions, that work with default logger/instrumenter (which is not exposed)
- package is configured with ENV vars, like `LOG_TAGS=app:program go run program.go`

Recommended way to use:
- import and use `"github.com/bsm/rucksack/log"` or `"github.com/bsm/rucksack/met"` in non-main packages
- import extension packages like `_ "github.com/bsm/rucksack/log/kafka"` only in main package


## Logging

[![GoDoc](https://godoc.org/github.com/bsm/rucksack/log?status.svg)](https://godoc.org/github.com/bsm/rucksack/log)

```
import (
	"github.com/bsm/rucksack/log"
	_ "github.com/bsm/rucksack/log/kafka"
)
```

Recommended ENV:
- `LOG_TAGS=foo:bar,baz:qux`
- `LOG_KAFKA_TOPIC=projectname` (with included [log/kafka](https://godoc.org/github.com/bsm/rucksack/log/kafka))
- `LOG_KAFKA_ADDRS=broker-1:9092,broker-2:9092` (with included [log/kafka](https://godoc.org/github.com/bsm/rucksack/log/kafka))

Optional ENV:
- `LOG_LEVEL=INFO`
- `LOG_KAFKA_LEVEL=INFO` (with included [log/kafka](https://godoc.org/github.com/bsm/rucksack/log/kafka))


## Metrics

[![GoDoc](https://godoc.org/github.com/bsm/rucksack/met?status.svg)](https://godoc.org/github.com/bsm/rucksack/met)

```
import (
	"github.com/bsm/rucksack/met"
	_ "github.com/bsm/rucksack/met/datadog"
	_ "github.com/bsm/rucksack/met/runtime"
)
```

Required/recommended ENV:
- `MET_NAME=projectname` (required)
- `MET_TAGS=foo:bar,baz:qux`
- `MET_DATADOG=datadog-token` (with included [met/datadog](https://godoc.org/github.com/bsm/rucksack/met/datadog); flush metrics to [datadog](https://www.datadoghq.com/))

Optional ENV:
- `HOST=hostname` (auto-detected)
- `PORT=8080` (optional)
- `MET_RUNTIME=mem,heap,gc` (with included [met/runtime](https://godoc.org/github.com/bsm/rucksack/met/runtime); which runtime metrics to report; used set is equivalent to `all`)

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
