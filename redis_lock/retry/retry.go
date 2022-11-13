package retry

import "time"

type Retry interface {
	Next() (bool, time.Duration)
}

type fixRetry struct {
	used int
	max  int
}

func (f fixRetry) Next() (bool, time.Duration) {
	f.used++
	return f.used <= f.max, 1
}
