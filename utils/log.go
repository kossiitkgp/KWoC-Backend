package utils

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func LogErr(r *http.Request, err error, errMsg string) {
	log.Err(err).Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		errMsg,
	)
}

func LogErrAndRespond(r *http.Request, w http.ResponseWriter, err error, errMsg string, statusCode int) {
	LogErr(r, err, errMsg)
	RespondWithHTTPMessage(r, w, statusCode, errMsg)
}

func LogInfo(r *http.Request, info string) {
	log.Info().Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		info,
	)
}

func LogWarn(r *http.Request, warning string) {
	log.Warn().Msgf(
		"%s %s: %s",
		r.Method,
		r.RequestURI,
		warning,
	)
}

func LogWarnAndRespond(r *http.Request, w http.ResponseWriter, warning string, statusCode int) {
	LogWarn(r, warning)
	RespondWithHTTPMessage(r, w, statusCode, warning)
}
