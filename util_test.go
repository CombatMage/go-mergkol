package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExist(t *testing.T) {
	// action
	result := fileExist("testdata")
	// verify
	assert.True(t, result)
}
