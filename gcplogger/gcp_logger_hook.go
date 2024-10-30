package gcplogger

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

// refer https://cloud.google.com/go/docs/reference/cloud.google.com/go/logging/latest
type gcpLoggingHook struct {
	client *logging.Client
	logger *logging.Logger
}

func newGcpLoggingHook(ctx context.Context, logID string, projectId string, tokenSource oauth2.TokenSource) (*gcpLoggingHook, error) {
	client, err := logging.NewClient(ctx, projectId, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}

	hook := &gcpLoggingHook{
		client: client,
		logger: client.Logger(logID),
	}

	hook.client.OnError = hook.onError

	return hook, nil
}

func (hook *gcpLoggingHook) onError(err error) {
	// TODO change it to logger and test this
	fmt.Println(fmt.Sprintf("Error detected from stackdriver: %+v", err))
}

func (hook *gcpLoggingHook) Close() error {
	return hook.client.Close()
}

func (hook *gcpLoggingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *gcpLoggingHook) Fire(entry *logrus.Entry) error {
	payload := map[string]interface{}{
		"message": entry.Message,
		"fields":  entry.Data,
	}

	if entry.HasCaller() {
		payload["reportLocation"] = map[string]interface{}{
			"filePath":     entry.Caller.File,
			"functionName": entry.Caller.Function,
			"lineNumber":   entry.Caller.Line,
		}
	}

	if errValue, ok := entry.Data[logrus.ErrorKey]; ok {
		if err, isErr := errValue.(error); isErr {
			payload["error"] = err.Error()
		}
	}

	severity := getSeverity(entry.Level)
	hook.logger.Log(logging.Entry{
		Payload:  payload,
		Severity: severity,
	})

	if severity >= logging.Error {
		err := formatError(entry)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func getSeverity(level logrus.Level) logging.Severity {
	switch level {
	case logrus.DebugLevel:
		return logging.Debug
	case logrus.InfoLevel:
		return logging.Info
	case logrus.WarnLevel:
		return logging.Warning
	case logrus.ErrorLevel:
		return logging.Error
	case logrus.FatalLevel, logrus.PanicLevel:
		return logging.Critical
	default:
		return logging.Default
	}
}

// TODO not a nice way to fmt
func formatError(entry *logrus.Entry) error {
	return fmt.Errorf("error: %v, Code: %v, Description: %v, Message: %s, Env: %v",
		entry.Data["error"], entry.Data["ErrCode"], entry.Data["ErrDescription"],
		entry.Message, entry.Data["env"])
}
