package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

const version = "0.1.0"

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
	if os.Args[1] == "-V" || os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || isValidPath(os.Args[1]) {
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
	fmt.Println(` codemeld "<file1> <file2> <file3> ..."`)
	fmt.Println(` codemeld -V 输出版本号`)
	fmt.Println("\n描述 | Description:")
	fmt.Println(" 将多个文件的内容以LLM更加容易理解的方式格式化并复制到剪贴板。")
	fmt.Println(" Format and copy the contents of multiple files to the clipboard in a way that is easier for LLMs to understand.")
	fmt.Println("\n示例 | Example:")
	fmt.Println(` codemeld "./file1.txt /home/user/file2.go ~/file3.md"`)
}

func isValidPath(path string) bool {
	// 简单路径检查（根据需求可以扩展更复杂的检查）
	// 在这个例子中，我们只做了一个简单的检查：路径是否包含非法字符
	// 更复杂的检查可以使用正则表达式
	re := regexp.MustCompile(`^[\w\-/\.]+$`)
	return re.MatchString(path)
}
