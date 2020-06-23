package main

import (
	"path/filepath"
	"testing"

	"github.com/EricNeid/go-mergkol/internal/verify"
)

func TestRemoveDuplicatesUnordered(t *testing.T) {
	// arrange
	data := []string{"foo", "bar", "foo"}
	// action
	result := removeDuplicatesUnordered(data)
	// verify
	verify.Equals(t, 2, len(result))
	verify.Contains(t, result, "bar")
}

func TestGetPathsForSourceCodeFiles(t *testing.T) {
	// should return only .kt files
	// action
	result, err := getPathsForSourceCodeFiles("testdata", ".kt", true)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 2, len(result))
	verify.Equals(t, "testdata/Main.kt", filepath.ToSlash(result[0]))
	verify.Equals(t, "testdata/test_pkg/Foo.kt", filepath.ToSlash(result[1]))

	// should return all files
	// action
	result, err = getPathsForSourceCodeFiles("testdata", "*", false)
	// verify
	verify.Ok(t, err)
	verify.Equals(t, 5, len(result))
	verify.Equals(t, "testdata/Main.kt", filepath.ToSlash(result[0])) // order according to file system
	verify.Equals(t, "testdata/MainTest.kt", filepath.ToSlash(result[1]))
	verify.Equals(t, "testdata/NoKotlin.java", filepath.ToSlash(result[2]))
	verify.Equals(t, "testdata/TestMain.kt", filepath.ToSlash(result[3]))
	verify.Equals(t, "testdata/test_pkg/Foo.kt", filepath.ToSlash(result[4]))
}

func TestNewSourceCodeFile(t *testing.T) {
	// action
	result, err := newSourceCodeFile("testdata/test_pkg/Foo.kt")
	// verify
	verify.Ok(t, err)
	verify.Equals(t, "testdata/test_pkg/Foo.kt", result.name)
	verify.Equals(t, "test_pkg", result.pkg)
	verify.Equals(t, 1, len(result.imports))
	verify.Equals(t, "class Foo(val bar: Int)", result.code[0])
}

func TestMergeSourceCodeFiles(t *testing.T) {
	// arrange
	file1, _ := newSourceCodeFile("testdata/test_pkg/Foo.kt")
	file2, _ := newSourceCodeFile("testdata/Main.kt")
	// action
	result := mergeSourceCodeFiles([]sourceCodeFile{file1, file2})
	// verify
	verify.Equals(t, 2, len(result.imports))
	verify.Contains(t, result.imports, "import java.util.*")
	verify.Contains(t, result.imports, "import java.lang.Math.abs")
	verify.Equals(t, "class Foo(val bar: Int)", result.code[0])
}

func TestFileExist_shouldReturnTrue(t *testing.T) {
	// action
	result := fileExist("testdata")
	// verify
	verify.Assert(t, result == true, "testdata not found")
}

func TestFileExist_shouldReturnFalse(t *testing.T) {
	// action
	result := fileExist("no_folder")
	// verify
	verify.Assert(t, result == false, "no_folder should not be found")
}
