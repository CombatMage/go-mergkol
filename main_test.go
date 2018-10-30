package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicatesUnordered(t *testing.T) {
	// arrange
	data := []string{"foo", "bar", "foo"}
	// action
	result := removeDuplicatesUnordered(data)
	// verify
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, "foo")
	assert.Contains(t, result, "bar")
}

func TestGetPathsForSourceCodeFiles(t *testing.T) {
	// should return only .kt files
	// action
	result, err := getPathsForSourceCodeFiles("testdata", ".kt", true)
	// verify
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "testdata/Main.kt", filepath.ToSlash(result[0]))
	assert.Equal(t, "testdata/test_pkg/Foo.kt", filepath.ToSlash(result[1]))

	// should return all files
	// action
	result, err = getPathsForSourceCodeFiles("testdata", "*", false)
	// verify
	assert.Nil(t, err)
	assert.Equal(t, 5, len(result))
	assert.Equal(t, "testdata/Main.kt", filepath.ToSlash(result[0])) // order according to file system
	assert.Equal(t, "testdata/MainTest.kt", filepath.ToSlash(result[1]))
	assert.Equal(t, "testdata/NoKotlin.java", filepath.ToSlash(result[2]))
	assert.Equal(t, "testdata/TestMain.kt", filepath.ToSlash(result[3]))
	assert.Equal(t, "testdata/test_pkg/Foo.kt", filepath.ToSlash(result[4]))
}

func TestNewSourceCodeFile(t *testing.T) {
	// action
	result, err := newSourceCodeFile("testdata/test_pkg/Foo.kt")
	// verify
	assert.Nil(t, err)
	assert.Equal(t, "testdata/test_pkg/Foo.kt", result.name)
	assert.Equal(t, "test_pkg", result.pkg)
	assert.Equal(t, 1, len(result.imports))
	assert.Equal(t, "class Foo(val bar: Int)", result.code[0])
}

func TestMergeSourceCodeFiles(t *testing.T) {
	// arrange
	file1, _ := newSourceCodeFile("testdata/test_pkg/Foo.kt")
	file2, _ := newSourceCodeFile("testdata/Main.kt")
	// action
	result := mergeSourceCodeFiles([]sourceCodeFile{file1, file2})
	// verify
	assert.Equal(t, 2, len(result.imports))
	assert.Contains(t, result.imports, "import java.util.*", "import java.lang.Math.abs")
	assert.Equal(t, "class Foo(val bar: Int)", result.code[0])
}
