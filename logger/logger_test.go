package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggerLogLevel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		expected string
	}{
		{
			logLevel: "info",
			expected: "info",
		},
		{
			logLevel: "warn",
			expected: "warn",
		},
		{
			logLevel: "debug",
			expected: "debug",
		},
		{
			logLevel: "error",
			expected: "error",
		},
		{
			logLevel: "fatal",
			expected: "fatal",
		},
		{
			logLevel: "trace",
			expected: "trace",
		},
		{
			logLevel: "panic",
			expected: "panic",
		},
		{
			logLevel: "plop",
			expected: "info",
		},
	}

	for _, tc := range tests {
		os.Setenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL", tc.logLevel)
		SetLoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(tc.expected, z)
		os.Unsetenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL")
	}
}
