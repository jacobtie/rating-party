package logger

import (
	"context"

	"github.com/jacobtie/rating-party/server/internal/platform/config"
	"github.com/jacobtie/rating-party/server/internal/platform/contextvalue"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger *zerolog.Logger

func InitLogger(environment config.ENV, instance, version, revision, buildDate string) {
	log.Level(zerolog.InfoLevel)
	if environment == config.ENV_LOCAL {
		log.Level(zerolog.DebugLevel)
	}
	l := log.With().
		Str("instance", instance).
		Str("version", version).
		Str("revision", revision).
		Str("buildDate", buildDate).
		Logger()
	globalLogger = &l
}

func Get() *zerolog.Logger {
	return globalLogger
}

func GetCtx(ctx context.Context) *zerolog.Logger {
	v, ok := ctx.Value(contextvalue.KeyValues).(*contextvalue.Values)
	if !ok {
		return Get()
	}
	l := globalLogger.With().
		Str("requestId", v.RequestID).
		Str("requestStart", v.RequestStart.String()).
		Str("method", v.Method).
		Str("path", v.Path).
		Str("ip", v.IP).
		Str("host", v.Host).
		Str("referrer", v.Referrer).
		Logger()
	return &l
}
