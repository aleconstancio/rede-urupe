package timeutil

import (
	"time"
	"github.com/aleconstancio/minos/timeutil"
)

var BrasiliaLocation = timeutil.BrasiliaLocation
var Now = timeutil.Now
var InBrasilia = timeutil.InBrasilia

func Today() string {
	return time.Now().In(BrasiliaLocation).Format("2006-01-02")
}
