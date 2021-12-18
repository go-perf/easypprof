package easypprof_test

import (
	"context"
	"flag"
	"time"

	"github.com/go-perf/easypprof"
)

func ExampleEasyPprof() {
	p, err := easypprof.NewProfiler(easypprof.Config{
		ProfileMode:   "mutex",
		FilePrefix:    "abc",
		UseTextFormat: false,
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := p.Run(ctx); err != nil {
		panic(err)
	}

	// Output:
}

func ExampleEasierPprof() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode:   "heap",
		FilePrefix:    "xyz",
		UseTextFormat: false,
	})

	// Output:
}

func ExampleCpuAsDefault() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{})

	// Output:
}

func ExampleCpuProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.CpuMode,
	})

	// Output:
}

func ExampleTraceProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.TraceMode,
	})

	// Output:
}

func ExampleHeapProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.HeapMode,
	})

	// Output:
}

func ExampleAllocsProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.AllocsMode,
	})

	// Output:
}

func ExampleMutexProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.MutexMode,
	})

	// Output:
}

func ExampleBlockProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.BlockMode,
	})

	// Output:
}

func ExampleThreadCreateProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.ThreadCreateMode,
	})

	// Output:
}

func ExampleGoroutineProfile() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: easypprof.GoroutineMode,
	})

	// Output:
}

func ExampleProfilePath() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go easypprof.Run(ctx, easypprof.Config{})

	// Output:
}

func ExampleWithFlags() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mode := flag.String("profile.mode", "", "enable profiling mode, one of: cpu, goroutine, threadcreate, heap, allocs, block, mutex.")
	disable := flag.Bool("profile.disable", false, "disable profiling? default false")
	flag.Parse()

	go easypprof.Run(ctx, easypprof.Config{
		ProfileMode: *mode,
		Disable:     *disable,
	})

	// Output:
}
