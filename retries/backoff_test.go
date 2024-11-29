package retries

import (
	"errors"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	t.Run("completes on success", func(t *testing.T) {
		sleepCalls := []time.Duration{}
		sleep = func(d time.Duration) {
			sleepCalls = append(sleepCalls, d)
		}

		var want int
		f := func() (Done, error) {
			want++
			return true, nil
		}

		err := Backoff(3, 1, f)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if want != 1 {
			t.Errorf("expected 1, got %v", want)
		}

		if len(sleepCalls) != 0 {
			t.Errorf("expected 0, got %v", len(sleepCalls))
		}
	})

	t.Run("returns errs if func errs", func(t *testing.T) {
		sleepCalls := []time.Duration{}
		sleep = func(d time.Duration) {
			sleepCalls = append(sleepCalls, d)
		}

		f := func() (Done, error) {
			return false, errors.New("some error")
		}

		err := Backoff(3, 1, f)

		if err.Error() != "some error" {
			t.Errorf("expected some error, got %v", err)
		}

		if len(sleepCalls) != 0 {
			t.Errorf("expected 0, got %v", len(sleepCalls))
		}
	})

	t.Run("retries until completion", func(t *testing.T) {
		sleepCalls := []time.Duration{}
		sleep = func(d time.Duration) {
			sleepCalls = append(sleepCalls, d)
		}

		var want int
		f := func() (Done, error) {
			want++
			if want < 3 {
				return false, nil
			}

			return true, nil
		}

		err := Backoff(3, 1, f)

		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		if want != 3 {
			t.Errorf("expected 3, got %v", want)
		}

		if len(sleepCalls) != 2 {
			t.Errorf("expected 2, got %v", len(sleepCalls))
		}

		wantSleepCalls := []time.Duration{
			0 * time.Second,
			1 * time.Second,
		}

		for i, v := range sleepCalls {
			if v != wantSleepCalls[i] {
				t.Errorf("expected %v, got %v", wantSleepCalls[i], v)
			}
		}
	})

	t.Run("returns exhausted error", func(t *testing.T) {
		sleepCalls := []time.Duration{}
		sleep = func(d time.Duration) {
			sleepCalls = append(sleepCalls, d)
		}

		f := func() (Done, error) {
			return false, nil
		}

		err := Backoff(3, 2, f)

		if err.Error() != ErrBackoffExhausted {
			t.Errorf("expected %v, got %v", ErrBackoffExhausted, err)
		}

		if len(sleepCalls) != 3 {
			t.Errorf("expected 3, got %v", len(sleepCalls))
		}

		wantSleepCalls := []time.Duration{
			0 * time.Second,
			2 * time.Second,
			4 * time.Second,
		}

		for i, v := range sleepCalls {
			if v != wantSleepCalls[i] {
				t.Errorf("expected %v, got %v", wantSleepCalls[i], v)
			}
		}
	})
}
