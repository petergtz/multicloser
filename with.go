package closers

import "io"

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
