// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package log

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger set the default logger to development mode
var Logger *zap.SugaredLogger

func init() {
	ConfigureLogger("dev")
}

func ConfigureLogger(logType string) {
	var cfg zap.Config
	if logType == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.LevelKey = "severity"
		cfg.EncoderConfig.MessageKey = "message"
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalln("createLogger", err)
	}
	Logger = logger.Sugar()
}

var CliLogger = createCliLogger()

func createCliLogger() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.LevelKey = ""
	cfg.DisableCaller = true
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalln("createLogger", err)
	}

	return logger.Sugar()
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, middleware.RequestIDKey, id)
}

func RequestIDLogger(r *http.Request) *zap.SugaredLogger {
	return ContextLogger(r.Context())
}

func ContextLogger(ctx context.Context) *zap.SugaredLogger {
	proposedLogger := Logger
	if ctx != nil {
		if ctxRequestID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
			proposedLogger = proposedLogger.With(zap.String("requestID", ctxRequestID))
		}
	}

	return proposedLogger
}
