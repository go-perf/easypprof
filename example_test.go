package easypprof_test

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/go-perf/easypprof"
)

func ExampleEasyPprof() {
	cfg := easypprof.Config{
		Disable:              os.Getenv("NO_EASY_PPROF") == "true",
		ProfileMode:          easypprof.CpuMode,
		OutputDir:            path.Join(".", "test_pprof"),
		FilePrefix:           "my-app",
		UseTextFormat:        false,
		MutexProfileFraction: 12,
		BlockProfileRate:     100_000,
		FgprofFormat:         "pprof",
	}
	defer easypprof.Start(cfg).Stop()

	// Output:
}

func ExampleEasyPprofWithPause() {
	// default config is used
	p := easypprof.Start(easypprof.Config{})

	defer func() {
		time.Sleep(time.Second)
		// and then stop the profiling
		p.Stop()
	}()

	// Output:
}

func ExampleWithEnv() {
	disable := os.Getenv("EASY_PPROF_DISABLE") == "true"
	mode := os.Getenv("EASY_PPROF_MODE")

	defer easypprof.Start(easypprof.Config{
		Disable:     disable,
		ProfileMode: mode,
		OutputDir:   path.Join(".", "test_pprof"),
		FilePrefix:  "env",
	}).Stop()

	// Output:
}

func ExampleWithFlags() {
	mode := flag.String("profile.mode", "", "enable profiling mode, one of: cpu, goroutine, threadcreate, heap, allocs, block, mutex.")
	disable := flag.Bool("profile.disable", false, "disable profiling? default false")
	flag.Parse()

	defer easypprof.Start(easypprof.Config{
		Disable:     *disable,
		ProfileMode: *mode,
		OutputDir:   path.Join(".", "test_pprof"),
		FilePrefix:  "flag",
	}).Stop()

	// Output:
}
