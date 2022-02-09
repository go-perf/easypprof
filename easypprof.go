package easypprof

import (
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

// Start profiling based on a config.
func Start(cfg Config) *Profiler {
	p, err := newProfiler(cfg)
	if err != nil {
		panic(fmt.Errorf("easypprof: %w", err))
	}
	if err := p.run(); err != nil {
		panic(fmt.Errorf("easypprof: %w", err))
	}
	return p
}

// Profiler for the Go programs.
type Profiler struct {
	cfg Config

	mode       string
	output     io.WriteCloser
	fgprofStop func() error // a bit hacky but simple

	memProfileRate       int
	blockProfileRate     int
	mutexProfileFraction int
	fgprofFormat         fgprof.Format
}

// Config of the profiler.
type Config struct {
	// Disable the profiler. Easy to set as a command-line flag.
	// Default is false.
	Disable bool

	// Mode is one of cpu, goroutine, threadcreate, heap, allocs, block, mutex.
	// Default is empty string and is treated as cpu.
	Mode string

	// OutputDir where profile must be saved.
	// Default is empty string and is treated as ".".
	OutputDir string

	// FilePrefix for the output file.
	// If prefix is specified an additional '_' will be added to the prefix.
	FilePrefix string

	// UseTextFormat of the resulting pprof file.
	// Default is false.
	UseTextFormat bool

	// FgprofFormat is a format of the output file.
	// Default is fgprof.FormatPprof.
	FgprofFormat string

	// MutexProfileFraction represents the fraction of mutex contention events.
	// Default is 10.
	MutexProfileFraction int

	// BlockProfileRate represents the fraction of goroutine blocking events.
	// Default is 10_000.
	BlockProfileRate int
}

// newProfiler creates new profiler based on a config.
func newProfiler(cfg Config) (*Profiler, error) {
	if cfg.Mode == "" {
		cfg.Mode = CpuMode
	}

	switch cfg.Mode {
	case CpuMode, TraceMode, HeapMode, AllocsMode, MutexMode,
		BlockMode, ThreadCreateMode, GoroutineMode, FgprofMode:
		// pass
	default:
		return nil, fmt.Errorf("unknown profile mode: %s", cfg.Mode)
	}

	if cfg.OutputDir == "" {
		cfg.OutputDir = "."
	}
	if cfg.FgprofFormat == "" {
		cfg.FgprofFormat = string(fgprof.FormatPprof)
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
	filename := fmt.Sprintf("%s%s_%s.pprof", prefix, cfg.Mode, now)

	if err := os.MkdirAll(cfg.OutputDir, 0777); err != nil {
		return nil, err
	}
	output, err := os.Create(path.Join(cfg.OutputDir, filename))
	if err != nil {
		return nil, err
	}

	p := &Profiler{
		cfg:    cfg,
		mode:   cfg.Mode,
		output: output,
	}
	return p, nil
}

// Run the profiler.
func (p *Profiler) run() error {
	switch p.mode {
	case CpuMode:
		return pprof.StartCPUProfile(p.output)
	case TraceMode:
		return trace.Start(p.output)
	case MutexMode:
		runtime.SetMutexProfileFraction(p.mutexProfileFraction)
	case BlockMode:
		runtime.SetBlockProfileRate(p.blockProfileRate)
	case FgprofMode:
		p.fgprofStop = fgprof.Start(p.output, p.fgprofFormat)
	case HeapMode, AllocsMode, ThreadCreateMode, GoroutineMode:
		// pass
	}
	return nil
}

// Stop the profiler and save the result.
func (p *Profiler) Stop() error {
	switch p.mode {
	case CpuMode:
		pprof.StopCPUProfile()
	case TraceMode:
		trace.Stop()
	case MutexMode:
		runtime.SetMutexProfileFraction(0)
	case BlockMode:
		runtime.SetBlockProfileRate(0)
	case FgprofMode:
		if err := p.fgprofStop(); err != nil {
			return err
		}
	}

	switch p.mode {
	case CpuMode, TraceMode, FgprofMode:
		// skip
	default:
		profile := pprof.Lookup(p.mode)
		if err := profile.WriteTo(p.output, bool2int(p.cfg.UseTextFormat)); err != nil {
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
