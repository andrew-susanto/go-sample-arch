package log

import (
	// golang package
	"errors"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
)

func TestLog_InitLogger(t *testing.T) {
	InitLogger()
	assert.True(t, true)
}

func TestLog_Error(t *testing.T) {
	// test log error before init logger
	zerologInstance = nil
	Error(errors.New("testing"), map[string]interface{}{
		"error": "error",
	}, "something got error")

	assert.True(t, true)

	// test log error after init logger
	InitLogger()
	Error(errors.New("testing"), map[string]interface{}{
		"error": "error",
	}, "something got error")

	assert.True(t, true)
}

func TestLog_Info(t *testing.T) {
	// test log info before init logger
	zerologInstance = nil
	Info(map[string]interface{}{
		"error": "error",
	}, "something got error")

	assert.True(t, true)

	// test log info after init logger
	InitLogger()
	Info(map[string]interface{}{
		"error": "error",
	}, "something got error")

	assert.True(t, true)
}
