package logger

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level         string
		expectedError bool
	}{
		{"debug", false},
		{"info", false},
		{"warn", false},
		{"error", false},
		{"panic", false},
		{"fatal", false},
		{"invalid", true},
	}

	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			logger, err := GetLogger(test.level)
			if test.expectedError {
				require.Error(t, err)
				require.Nil(t, logger)
			} else {
				require.Nil(t, err)
				require.NotNil(t, logger)
			}
		})
	}
}

func TestLoggerWithCustomContext(t *testing.T) {
	logger, err := GetLogger("debug")
	require.NoError(t, err)
	ctx := ContextWithLogger(context.Background(), logger)
	require.NotNil(t, ctx.Value("logger"), "Expected logger in context")

	retrievedLogger, _ := GetLoggerFromContext(ctx)
	require.NotNil(t, retrievedLogger, "Expected logger retrieved from context")
	require.Equal(t, logger, retrievedLogger, "Retrieved logger does not match expected logger")
}

func TestLoggerFromContextWithoutLogger(t *testing.T) {
	ctx := context.Background()
	logger, err := GetLoggerFromContext(ctx)
	require.NoError(t, err)
	require.NotNil(t, logger, "Expected a new logger when not present in context")
}

func TestLogOutputToFileAndStdout(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	logger, err := GetLogger("debug")
	require.NoError(t, err)
	defer func() {
		os.Stdout = oldStdout
	}()

	logger.Info("Test logger message", nil)
	err = w.Close()
	require.NoError(t, err)

	got := make([]byte, 100)
	_, err = r.Read(got)
	require.NoError(t, err)

	require.Contains(t, string(got), "Test logger message", "Expected log message not found in file output")
}
