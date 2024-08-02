package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func getCommonRoot(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	if len(paths) == 1 {
		return filepath.Dir(paths[0])
	}

	slashPaths := make([]string, len(paths))
	for i, path := range paths {
		slashPaths[i] = filepath.ToSlash(path)
	}

	parts := strings.Split(slashPaths[0], "/")

	for i := 1; i < len(slashPaths); i++ {
		other := strings.Split(slashPaths[i], "/")
		j := 0
		for j < len(parts) && j < len(other) && parts[j] == other[j] {
			j++
		}
		parts = parts[:j]
	}

	commonPath := strings.Join(parts, "/")
	return filepath.FromSlash(commonPath)
}

func readFileContent(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	return string(content)
}

func processFiles(filePaths []string) (string, string, []string) {
	rootPath := getCommonRoot(filePaths)
	var result []string
	var relativePaths []string

	for _, path := range filePaths {
		content := readFileContent(path)
		relativePath, _ := filepath.Rel(rootPath, path)
		relativePaths = append(relativePaths, relativePath)
		result = append(result, fmt.Sprintf("```%s\n%s\n```", relativePath, content))
	}

	return rootPath, strings.Join(result, "\n\n"), relativePaths
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}

	inputString := os.Args[1]
	filePaths := strings.Fields(inputString)

	rootPath, formattedContent, relativePaths := processFiles(filePaths)

	fmt.Printf("Root path: %s\n\n", rootPath)
	fmt.Printf("%d file(s):\n", len(relativePaths))
	for _, path := range relativePaths {
		fmt.Printf("  - %s\n", path)
	}
	err := clipboard.WriteAll(formattedContent)
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
	} else {
		fmt.Println("Formatted content copied to clipboard.")
	}
}

func printUsage() {
	fmt.Println("用法 | Usage:")
	fmt.Println(" codemeld <file1> <file2> <file3> ...")
	fmt.Println("\n描述 | Description:")
	fmt.Println(" 将多个文件的内容以LLM更加容易理解的方式格式化并复制到剪贴板。")
	fmt.Println(" Format and copy the contents of multiple files to the clipboard in a way that is easier for LLMs to understand.")
	fmt.Println("\n示例 | Example:")
	fmt.Println(" codemeld ./file1.txt /home/user/file2.go ~/file3.md")
}
