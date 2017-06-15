package closers

import "io"

// With runs the code in run and automatically closes closer afterwards.
// In case run returns a non-nil error With will return that error after closing closer.
// If additionally closer.Close() returns an error, that error is suppressed.
// Only when run returns nil, but closing closer returns non-nil, will With return the error
// returned from Closing.
// close Error Convert can be used to add context to the error returned from Close.
// If no context is needed the Identity function can be used.
func With(closer io.Closer, run func() error, closeErrorConvert func(error) error) error {
	e := run()
	if e != nil {
		closer.Close()
		return e
	}
	e = closer.Close()
	if e != nil {
		return closeErrorConvert(e)
	}
	return nil
}

func Identity(e error) error {
	return e
}
