package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultFileExtensionFilter = "*"
const defaultDir = "src"
const defaultOutputFile = "Merged.kt"

// sourceCodeFile represents a file of source code.
// It contains name and package (both optional) and
// a list of imports and lines of code (both can be empty).
type sourceCodeFile struct {
	name    string
	pkg     string
	imports []string
	code    []string
}

func (file sourceCodeFile) writeToFile(path string) error {
	var mergedFile string
	for _, line := range file.imports {
		mergedFile += line + "\r\n"
	}
	mergedFile += "\r\n"
	for _, line := range file.code {
		mergedFile += line + "\r\n"
	}
	return ioutil.WriteFile(path, []byte(mergedFile), 777)
}

func newSourceCodeFile(path string) (sourceCodeFile, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return sourceCodeFile{}, err
	}

	var imports []string
	var code []string
	var pkg string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "import") {
			imports = append(imports, line)
		} else if strings.HasPrefix(line, "package") {
			pkg = strings.SplitN(line, "package", 2)[1]
			pkg = strings.TrimSpace(pkg)
		} else if strings.TrimSpace(line) != "" {
			code = append(code, line)
		}
	}

	return sourceCodeFile{
			name:    file.Name(),
			pkg:     pkg,
			imports: imports,
			code:    code,
		},
		scanner.Err()
}

// getPathsForSourceCodeFiles returns slice, containing the paths to files with
// supplied extensions. If wildcard (*) is given, all files are returned. Each path starts with the supplied dir.
// If skipTestFiles is set, each file with test (ignoring case) in its name is ignored.
func getPathsForSourceCodeFiles(dir string, fileExtensionFilter string, skipTestFiles bool) (paths []string, err error) {
	err = filepath.Walk(
		dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fileName := info.Name()
				// check if file is test
				if skipTestFiles && strings.Contains(strings.ToLower(fileName), "test") {
					return nil
					// check for matching extension
				} else if fileExtensionFilter == "*" || strings.HasSuffix(fileName, fileExtensionFilter) {
					paths = append(paths, path)
				}
			}
			return nil
		})
	return paths, err
}

// mergeSourceCodeFiles merges all files in the given slice into a single sourceCodeFile.
// Imports are filtered for duplicates and the code lines are simply joined.
// Lokal package information is stripped.
func mergeSourceCodeFiles(files []sourceCodeFile) sourceCodeFile {
	var imports []string
	var code []string
	var pkgs []string

	// aggregate content
	for _, file := range files {
		pkgs = append(pkgs, file.pkg)
		imports = append(imports, file.imports...)
		code = append(code, file.code...)
	}
	imports = removeDuplicatesUnordered(imports)

	// remove local package imports
	for _, pkg := range pkgs {
		if pkg != "" {
			imports = removePackageImport(imports, pkg)
		}
	}

	return sourceCodeFile{
		imports: imports,
		code:    code,
	}
}

// removePackageImport removes all imports of package pkg from imports
func removePackageImport(imports []string, pkg string) []string {
	var stripedImports []string
	for _, importLine := range imports {
		importedPkg := strings.SplitN(importLine, "import", 2)[1]
		importedPkg = strings.TrimSpace(importedPkg)
		if !strings.HasPrefix(importedPkg, pkg) {
			stripedImports = append(stripedImports, importLine)
		}
	}
	return stripedImports
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}
	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}
	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

// mergeFilesInDir discovers all files in the given directory with the extension.
// It can skip each files with test (ignoring case) in its names.
func mergeFilesInDir(dir string, extension string, skipTestFiles bool) (mergedFile sourceCodeFile, err error) {
	paths, err := getPathsForSourceCodeFiles(dir, extension, skipTestFiles)
	if err != nil {
		fmt.Println("error while file discovery: " + err.Error())
		return sourceCodeFile{}, err
	}

	var files []sourceCodeFile
	for _, path := range paths {
		fmt.Println("\t" + path)
		file, err := newSourceCodeFile(path)
		if err != nil {
			fmt.Println("\terror while processing: " + err.Error())
		}
		files = append(files, file)
	}

	return mergeSourceCodeFiles(files), nil
}

type input struct {
	inputDir         string
	extensionsFilter string
	outputFile       string
	skipTestFiles    bool
	showHelp         bool
}

func readCmdArguments() input {
	helpPtr := flag.Bool("h", false, "print help")

	inputDirPtr := flag.String("dir", defaultDir, "source code directory")
	extensionFilterPtr := flag.String("file", defaultFileExtensionFilter, "only process file with this extension")
	skipTestFiles := flag.Bool("t", false, "skip files with test in name (ignore case)")

	outputPtr := flag.String("o", defaultOutputFile, "write merged code into this file")
	flag.Parse()

	return input{
		showHelp:         *helpPtr,
		inputDir:         *inputDirPtr,
		extensionsFilter: *extensionFilterPtr,
		skipTestFiles:    *skipTestFiles,
		outputFile:       *outputPtr,
	}
}

func main() {
	input := readCmdArguments()

	if input.showHelp {
		fmt.Println("Usage: ")
		fmt.Println("go-mergkol.exe -dir <source code dir> -o <output file>")
		return
	}
	fmt.Println("merging files in: " + input.inputDir)
	fmt.Println("reading files with extension: " + input.extensionsFilter)
	fmt.Println("skipping test files: " + strconv.FormatBool(input.skipTestFiles))

	if !fileExist(input.inputDir) {
		fmt.Println("input directory not found: " + input.inputDir)
		return
	}

	mergedFile, err := mergeFilesInDir(input.inputDir, input.extensionsFilter, input.skipTestFiles)

	fmt.Println("Write output to: " + input.outputFile)
	err = mergedFile.writeToFile(input.outputFile)
	if err != nil {
		fmt.Println("error while writing result: " + err.Error())
	}
}
