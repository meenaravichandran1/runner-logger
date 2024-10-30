package gcplogger

import (
	"context"
	"fmt"
	"github.com/harness/runner/delegateshell/client"
	"github.com/sirupsen/logrus"
	"github/meenaravichandran1/runner-logger/logger"
	"path"
	"runtime"
	"strconv"
	"time"
)

// refer https://cloud.google.com/go/docs/reference/cloud.google.com/go/logging/latest
const (
	logFileName = "runner.log"
)

type GCPLogger struct {
	gcpLoggingHook *gcpLoggingHook
	ManagerClient  *client.ManagerClient
}

func NewGCPLogger(managerClient *client.ManagerClient) *GCPLogger {
	return &GCPLogger{
		ManagerClient: managerClient,
	}
}

func (gcpLogger *GCPLogger) StartGcpLogger(ctx context.Context) (bool, error) {
	tokenManager, err := NewTokenManager(ctx, gcpLogger.ManagerClient)
	if err != nil {
		return false, fmt.Errorf("failed to initialize token provider: %w", err)
	}

	gcpLogger.gcpLoggingHook, err = newGcpLoggingHook(ctx, logFileName, tokenManager.ProjectID, tokenManager)
	if err != nil {
		return false, fmt.Errorf("failed to create stack driver hook: %w", err)
	}

	gcpRunnerLogger := logger.CreateNewLogger()
	gcpRunnerLogger.LogrusLogger.AddHook(gcpLogger.gcpLoggingHook)
	gcpRunnerLogger.LogrusLogger.SetReportCaller(true)
	gcpRunnerLogger.LogrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	})

	logger.ChangeLogger(gcpRunnerLogger)

	return true, nil
}

func (gcpLogger *GCPLogger) StopGcpLogger() (bool, error) {
	err := gcpLogger.gcpLoggingHook.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}
