package timex

import (
	"math"
	"time"
)

var UTCTime = func() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.00Z")
}

func TimeWithinRange(timestamp1, timestamp2 int64, tolerance time.Duration) bool {
	diff := math.Abs(float64(timestamp1 - timestamp2))
	diffDuration := time.Duration(diff)
	return diffDuration <= tolerance
}
