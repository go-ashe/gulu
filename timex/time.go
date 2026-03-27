package timex

import "time"

var UTCTime = func() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.00Z")
}
