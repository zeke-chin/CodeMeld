package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const version = "0.2.1"

var args AppArgs

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

type AppArgs struct {
	Files         []string
	Reg           []string
	IReg          []string
	IncludeHidden bool
}

func createApp() *cli.App {
	var filesInput, regInput, iregInput string
	var includeHidden bool

	return &cli.App{
		Name:  "codemeld",
		Usage: "一个将code内容格式化成LLM更友好的格式的剪切板命令行工具",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "files",
				Aliases:     []string{"f"},
				Usage:       "文件、文件夹目录 可以使用空格或者换行隔开\teg: -f \"/path/to/directory file3.txt\" ",
				Destination: &filesInput,
			},
			&cli.StringFlag{
				Name:        "reg",
				Aliases:     []string{"r"},
				Usage:       "匹配文件后缀 使用空格或者换行隔开\teg: -r \".go .py\"",
				Destination: &regInput,
			},
			&cli.StringFlag{
				Name:        "ireg",
				Aliases:     []string{"ir"},
				Usage:       "忽略的文件后缀 使用空格或者换行隔开\teg: -ir \".log .tmp\"",
				Destination: &iregInput,
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v", "V"},
				Usage:   "CodeMeld 版本",
			},
			&cli.BoolFlag{
				Name:        "include-hidden",
				Aliases:     []string{"ih"},
				Usage:       "是否查看隐藏文件夹下的文件",
				Destination: &includeHidden,
			},
		},
		Before: func(c *cli.Context) error {
			args = AppArgs{
				Files:         strings.Fields(filesInput),
				Reg:           strings.Fields(regInput),
				IReg:          strings.Fields(iregInput),
				IncludeHidden: includeHidden,
			}
			return nil
		},
		Action: handleAction,
	}
}

func handleAction(c *cli.Context) error {
	if c.Bool("version") {
		fmt.Printf("codemeld version: %v\n", version)
		os.Exit(0)
	}

	if len(args.Files) == 0 {
		_ = cli.ShowAppHelp(c)
		os.Exit(0)
	}

	return nil
}

func filtrationFiles(args AppArgs) ([]string, error) {
	var result []string

	for _, path := range args.Files {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if !args.IncludeHidden && strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}

			if shouldIncludeFile(filePath, args) {
				result = append(result, filePath)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func shouldIncludeFile(filePath string, args AppArgs) bool {
	ext := filepath.Ext(filePath)

	// Check if the file should be ignored
	for _, iregExt := range args.IReg {
		if strings.EqualFold(ext, iregExt) {
			return false
		}
	}

	// If no reg specified, include all files
	if len(args.Reg) == 0 {
		return true
	}

	// Check if the file matches the reg
	for _, regExt := range args.Reg {
		if strings.EqualFold(ext, regExt) {
			return true
		}
	}

	return false
}

func main() {
	app := createApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	filePaths, err := filtrationFiles(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	rootPath, formattedContent, relativePaths := processFiles(filePaths)

	fmt.Printf("Root path: %s\n\n", rootPath)
	fmt.Printf("%d file(s):\n", len(relativePaths))
	for _, path := range relativePaths {
		fmt.Printf("  - %s\n", path)
	}
	err = clipboard.WriteAll(formattedContent + "\n\n")
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
	} else {
		fmt.Println("Formatted content copied to clipboard.")
	}
}
