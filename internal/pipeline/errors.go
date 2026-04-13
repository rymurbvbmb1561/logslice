package pipeline

import "errors"

// Sentinel errors returned by Run.
var (
	// ErrNoReader is returned when Config.Reader is nil.
	ErrNoReader = errors.New("pipeline: Reader must not be nil")

	// ErrNoWriter is returned when Config.Writer is nil.
	ErrNoWriter = errors.New("pipeline: Writer must not be nil")

	// ErrNoStats is returned when Config.Stats is nil.
	ErrNoStats = errors.New("pipeline: Stats must not be nil")
)

// validate checks that the required Config fields are populated.
func validate(cfg Config) error {
	if cfg.Reader == nil {
		return ErrNoReader
	}
	if cfg.Writer == nil {
		return ErrNoWriter
	}
	if cfg.Stats == nil {
		return ErrNoStats
	}
	return nil
}
