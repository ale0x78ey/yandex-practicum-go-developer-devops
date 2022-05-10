package agent

import "time"

func minDuration(d1, d2 time.Duration) time.Duration {
	if d1 <= d2 {
		return d1
	}
	return d2
}
