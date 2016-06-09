# Rucksack

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
- `LOG_TAGS=dc:eu1,region:eu`
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
- `MET_TAGS=dc:eu1,region:eu`
- `MET_DATADOG=datadog-token` (with included [met/datadog](https://godoc.org/github.com/bsm/rucksack/met/datadog); flush metrics to [datadog](https://www.datadoghq.com/))

Optional ENV:
- `HOST=hostname` (auto-detected)
- `PORT=8080` (optional)
- `MET_RUNTIME=mem,heap,gc` (with included [met/runtime](https://godoc.org/github.com/bsm/rucksack/met/runtime); which runtime metrics to report; used set is equivalent to `all`)


## Licence

```
Copyright (c) 2016 Black Square Media

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
