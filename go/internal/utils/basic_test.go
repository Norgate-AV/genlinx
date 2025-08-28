package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicNormalizePath(t *testing.T) {
	result := NormalizePath("test/path")
	assert.NotEmpty(t, result)
}

func TestBasicFileExists(t *testing.T) {
	exists := FileExists("nonexistent")
	assert.False(t, exists)
}
