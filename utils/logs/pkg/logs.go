package logs

import (

	"os"

	kitlog "github.com/go-kit/kit/log"
)

var Logger kitlog.Logger

func init() {

	Logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	Logger = kitlog.With(Logger, "timestamp", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
}