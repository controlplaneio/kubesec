package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sublimino/kubesec/pkg/server"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "os"
	"strconv"
	"time"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   `http [port]`,
	Short: "Starts kubesec HTTP server on the specified port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("port is required")
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = args[0]
		}

		if _, err := strconv.Atoi(port); err != nil {
			port = args[0]
		}

		stopCh := server.SetupSignalHandler()
		jsonLogger, _ := NewJsonLogger("info")

		server.ListenAndServe(port, time.Minute, jsonLogger, stopCh)
		return nil
	},
}

// NewJsonLogger returns a zap sugared logger configured with json format recognized by Stackdriver
func NewJsonLogger(logLevel string) (*zap.SugaredLogger, error) {
  level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
  switch logLevel {
  case "debug":
    level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
  case "info":
    level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
  case "warn":
    level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
  case "error":
    level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
  case "fatal":
    level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
  case "panic":
    level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
  }

  zapEncoderConfig := zapcore.EncoderConfig{
    TimeKey:        "timestamp",
    LevelKey:       "severity",
    NameKey:        "logger",
    CallerKey:      "caller",
    MessageKey:     "message",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.LowercaseLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
  }

  zapConfig := zap.Config{
    Level:       level,
    Development: false,
    Sampling: &zap.SamplingConfig{
      Initial:    100,
      Thereafter: 100,
    },
    Encoding:         "json",
    EncoderConfig:    zapEncoderConfig,
    OutputPaths:      []string{"stderr"},
    ErrorOutputPaths: []string{"stderr"},
  }

  logger, err := zapConfig.Build()
  if err != nil {
    return nil, err
  }
  return logger.Sugar(), nil
}
