package retries

import (
	"errors"
	"time"
)

type Done = bool
type BackoffFunc func() (Done, error)

const ErrBackoffExhausted = "backoff retries exhausted"

var sleep = time.Sleep

func Backoff(maxRetries int, step int, f BackoffFunc) error {
	for i := 0; i < maxRetries; i++ {
		backoff := time.Duration(i*step) * time.Second
		done, err := f()

		if err != nil {
			return err
		}

		if done {
			return nil
		}

		if i < maxRetries {
			sleep(backoff)
		}
	}

	return errors.New(ErrBackoffExhausted)
}
