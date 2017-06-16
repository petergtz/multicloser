package multicloser

import (
	"io"
)

type MultiCloser struct {
	closers []io.Closer
}

// New creates a MultiCloser
func New(closers ...io.Closer) *MultiCloser {
	return &MultiCloser{closers}
}

// Add adds a new io.Closer to the MultiCloser. It's important to add a Closer to the MultiCloser as soon as it is available,
// to make sure it gets closed properly
func (mc *MultiCloser) Add(closer ...io.Closer) {
	mc.closers = append(mc.closers, closer...)
}

func (mc *MultiCloser) CloseAfter(block func() error, closeErrorConvert func(error) error) error {
	e := block()
	if e != nil {
		mc.Close()
		return e
	}
	e = mc.Close()
	if e != nil {
		if closeErrorConvert != nil {
			return closeErrorConvert(e)
		}
		return e
	}
	return nil

}

// Close closes all closers of MultiCloser in reverse order
func (mc *MultiCloser) Close() error {
	var err error
	numClosers := len(mc.closers)
	for i := range mc.closers {
		if closeErr := mc.closers[numClosers-i-1].Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}
