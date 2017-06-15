package closers

import (
	"io"
)

type multiCloser struct {
	closers []io.Closer
}

// Multi create a multi-Closer that will close all closers passed in as parameters
// in reverse order
func Multi(closers ...io.Closer) io.Closer {
	result := &multiCloser{make([]io.Closer, len(closers))}
	for i, closer := range closers {
		result.closers[len(closers)-i-1] = closer
	}
	return result
}

// Close closes all closers of the multi-Closer in reverse order
func (mc *multiCloser) Close() error {
	var err error
	for _, c := range mc.closers {
		if closeErr := c.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}
