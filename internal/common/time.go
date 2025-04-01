package common

import "time"

// Return current unix timestamp in seconds and UTC timezone.
func getStdTime() int64 {
	return time.Now().UTC().Unix()
}
