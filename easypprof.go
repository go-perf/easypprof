package easypprof

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

const (
	CpuMode          = "cpu"
	TraceMode        = "trace"
	HeapMode         = "heap"
	AllocsMode       = "allocs"
	MutexMode        = "mutex"
	BlockMode        = "block"
	ThreadCreateMode = "threadcreate"
	GoroutineMode    = "goroutine"
)

// Run ...
func Run(ctx context.Context, cfg Config) error {
	p, err := NewProfiler(cfg)
	if err != nil {
		return err
	}
	return p.Run(ctx)
}

// Profiler ...
type Profiler struct {
	profileMode   string
	output        io.WriteCloser
	useTextFormat bool
}

// Config ...
type Config struct {
	// Disable the profiler. Easy to set as a command-line flag. Default is false.
	Disable bool

	// ProfileMode is one of cpu, goroutine, threadcreate, heap, allocs, block, mutex.
	// Default is empty string and is treated as cpu.
	ProfileMode string

	// OutputDir where profile must be saved.
	// Default is empty string and is treated as ".".
	OutputDir string

	FilePrefix string

	// UseTextFormat of the resulting pprof file. Default is false.
	UseTextFormat bool

	// Profiles related parameters.
	MemProfileRate int
}

// NewProfiler ...
func NewProfiler(cfg Config) (*Profiler, error) {
	if cfg.ProfileMode == "" {
		cfg.ProfileMode = CpuMode
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "."
	}

	var prefix string
	if cfg.FilePrefix != "" {
		prefix = cfg.FilePrefix + "_"
	}
	now := time.Now().Format("20060102-15:04:05")
	filename := fmt.Sprintf("%s%s_%s.pprof", prefix, cfg.ProfileMode, now)

	output, err := os.Create(path.Join(cfg.OutputDir, filename))
	if err != nil {
		return nil, err
	}

	p := &Profiler{
		profileMode:   cfg.ProfileMode,
		output:        output,
		useTextFormat: cfg.UseTextFormat,
	}
	return p, nil
}

// Run ...
func (p *Profiler) Run(ctx context.Context) error {
	switch p.profileMode {
	case CpuMode:
		if err := pprof.StartCPUProfile(p.output); err != nil {
			return err
		}
	case TraceMode:
		if err := trace.Start(p.output); err != nil {
			return err
		}
	case MutexMode:
		runtime.SetMutexProfileFraction(1)
	case BlockMode:
		runtime.SetBlockProfileRate(1)
	default:
		// pass
	}

	<-ctx.Done()

	switch p.profileMode {
	case CpuMode:
		pprof.StopCPUProfile()
	case TraceMode:
		trace.Stop()
	case MutexMode:
		runtime.SetMutexProfileFraction(0)
	case BlockMode:
		runtime.SetBlockProfileRate(0)
	default:
		// pass
	}

	if p.profileMode != CpuMode && p.profileMode != TraceMode {
		profile := pprof.Lookup(p.profileMode)
		if err := profile.WriteTo(p.output, bool2int(p.useTextFormat)); err != nil {
			return err
		}
	}

	if err := p.output.Close(); err != nil {
		return err
	}
	return nil
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
