package utils

import "errors"

type ClosableIO interface {
	Close() error
}

func CloseAll(objects ...ClosableIO) error {
	var errs []error

	for _, object := range objects {
		if err := object.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
