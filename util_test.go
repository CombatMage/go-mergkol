package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExist_shouldReturnTrue(t *testing.T) {
	// action
	result := fileExist("testdata")
	// verify
	assert.True(t, result)
}

func TestFileExist_shouldReturnFalse(t *testing.T) {
	// action
	result := fileExist("no_folder")
	// verify
	assert.False(t, result)
}
