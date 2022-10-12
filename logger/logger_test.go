package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggerLogLevel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		expected string
		caller   bool
	}{
		{
			logLevel: "info",
			expected: "info",
			caller:   false,
		},
		{
			logLevel: "warn",
			expected: "warn",
			caller:   true,
		},
		{
			logLevel: "debug",
			expected: "debug",
			caller:   true,
		},
		{
			logLevel: "error",
			expected: "error",
			caller:   true,
		},
		{
			logLevel: "fatal",
			expected: "fatal",
			caller:   true,
		},
		{
			logLevel: "trace",
			expected: "trace",
			caller:   true,
		},
		{
			logLevel: "panic",
			expected: "panic",
			caller:   true,
		},
		{
			logLevel: "plop",
			expected: "info",
			caller:   false,
		},
	}

	for _, tc := range tests {
		if tc.caller {
			os.Setenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL_WITH_CALLER", "true")
		}
		os.Setenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL", tc.logLevel)
		SetLoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(tc.expected, z)
		os.Unsetenv("CYPRESS_PARALLEL_CLI_LOG_LEVEL")
	}
}

func TestLogger_info(t *testing.T) {
	SetLoggerLogLevel()
	log.Info().Msgf("Testing logger")
}
