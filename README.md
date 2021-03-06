# easypprof

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

Easy pprof library for Go.

## Note

Check [go-profiler-notes](https://github.com/DataDog/go-profiler-notes) for more information about profiling in Go.

## Features

* Simple API.
* Easy to integrate.
* Configurable.
* Supports [fgprof](https://github.com/felixge/fgprof).
* Improved version of [pkg/profile](https://github.com/pkg/profile).

## Install

Go version 1.17+

```
go get github.com/go-perf/easypprof
```

## Example

```go
func main() {
	cfg := easypprof.Config{
		Mode:       easypprof.CpuMode,
		OutputDir:  path.Join(".", "test_pprof"),
		FilePrefix: "my-app",
	}
	defer easypprof.Start(cfg).Stop()

	// your code
}
```

See examples for more: [examples_test.go](https://github.com/go-perf/easypprof/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/go-perf/easypprof/workflows/build/badge.svg
[build-url]: https://github.com/go-perf/easypprof/actions
[pkg-img]: https://pkg.go.dev/badge/go-perf/easypprof
[pkg-url]: https://pkg.go.dev/github.com/go-perf/easypprof
[reportcard-img]: https://goreportcard.com/badge/go-perf/easypprof
[reportcard-url]: https://goreportcard.com/report/go-perf/easypprof
[coverage-img]: https://codecov.io/gh/go-perf/easypprof/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/go-perf/easypprof
[version-img]: https://img.shields.io/github/v/release/go-perf/easypprof
[version-url]: https://github.com/go-perf/easypprof/releases
