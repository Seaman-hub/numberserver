// Package pools provides general purpose rate pool implementations.
package pools

import (
	"errors"
	"log"
)

var (
	// ErrLimitExhausted is returned by the Limiter in case the number of requests overflows the capacity of a Limiter.
	ErrLimitExhausted = errors.New("requests limit exhausted")

	// ErrRaceCondition is returned when there is a race condition while saving a state of a rate limiter.
	ErrRaceCondition = errors.New("race condition detected")
)

// Logger wraps the Log method for logging.
type Logger interface {
	// Log logs the given arguments.
	Log(v ...interface{})
}

// StdLogger implements the Logger interface.
type StdLogger struct{}

// NewStdLogger creates a new instance of StdLogger.
func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

// Log delegates the logging to the std logger.
func (l *StdLogger) Log(v ...interface{}) {
	log.Println(v...)
}
