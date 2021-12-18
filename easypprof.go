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

	"github.com/felixge/fgprof"
)

// Profile modes.
const (
	CpuMode          = "cpu"
	TraceMode        = "trace"
	HeapMode         = "heap"
	AllocsMode       = "allocs"
	MutexMode        = "mutex"
	BlockMode        = "block"
	ThreadCreateMode = "threadcreate"
	GoroutineMode    = "goroutine"
	FgprofMode       = "fgprof"
)

// Run a profiler based on a config.
func Run(ctx context.Context, cfg Config) error {
	p, err := NewProfiler(cfg)
	if err != nil {
		return err
	}
	return p.Run(ctx)
}

// Profiler for the Go programs.
type Profiler struct {
	profileMode   string
	output        io.WriteCloser
	useTextFormat bool

	memProfileRate       int
	blockProfileRate     int
	mutexProfileFraction int
	fgprofFormat         fgprof.Format
}

// Config of the profiler.
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

	// MutexProfileFraction represents the fraction of mutex contention events. Default is 10.
	MutexProfileFraction int

	// BlockProfileRate represents the fraction of goroutine blocking events. Default is 10_000.
	BlockProfileRate int

	// FgprofFormat is a format of the output file. Default is fgprof.FormatPprof.
	FgprofFormat fgprof.Format
}

// NewProfiler creates new Profile based on a config.
func NewProfiler(cfg Config) (*Profiler, error) {
	if cfg.ProfileMode == "" {
		cfg.ProfileMode = CpuMode
	}
	if cfg.OutputDir == "" {
		cfg.OutputDir = "."
	}
	if cfg.FgprofFormat == "" {
		cfg.FgprofFormat = fgprof.FormatPprof
	}
	if cfg.MutexProfileFraction == 0 {
		cfg.MutexProfileFraction = 10
	}
	if cfg.BlockProfileRate == 0 {
		cfg.BlockProfileRate = 10_000
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
		fgprofFormat:  cfg.FgprofFormat,
	}
	return p, nil
}

// Run the profiler.
func (p *Profiler) Run(ctx context.Context) error {
	// a bit hacky but simple
	var fgprofStop func() error

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
		runtime.SetMutexProfileFraction(p.mutexProfileFraction)
	case BlockMode:
		runtime.SetBlockProfileRate(p.blockProfileRate)
	case FgprofMode:
		fgprofStop = fgprof.Start(p.output, p.fgprofFormat)
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
	case FgprofMode:
		if err := fgprofStop(); err != nil {
			return err
		}
	default:
		// pass
	}

	switch p.profileMode {
	case CpuMode, TraceMode, FgprofMode:
		// skip
	default:
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
