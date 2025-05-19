package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const version = "0.2.2"
const maxCacheFiles = 50

var args AppArgs

// 生成8位随机ID
func generateID() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())[:8]
	}
	return hex.EncodeToString(b)
}

// 确保缓存目录存在
func ensureCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "codemeld")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", err
	}

	return cacheDir, nil
}

// 清理旧缓存文件
func cleanOldCacheFiles(cacheDir string) error {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	if len(files) <= maxCacheFiles {
		return nil
	}

	// 按修改时间排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	// 删除最旧的文件，直到文件数量不超过maxCacheFiles
	for i := 0; i < len(files)-maxCacheFiles; i++ {
		err := os.Remove(filepath.Join(cacheDir, files[i].Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

// 保存内容到缓存文件
func saveToCache(content string) (string, string, error) {
	cacheDir, err := ensureCacheDir()
	if err != nil {
		return "", "", fmt.Errorf("创建缓存目录失败: %w", err)
	}

	// 清理旧缓存文件
	err = cleanOldCacheFiles(cacheDir)
	if err != nil {
		fmt.Printf("警告: 清理旧缓存文件失败: %v\n", err)
		// 继续执行，不终止整个流程
	}

	// 生成ID和文件名
	id := generateID()
	fileName := fmt.Sprintf("%s_%s.md", time.Now().Format("20060102_150405"), id)
	filePath := filepath.Join(cacheDir, fileName)

	// 分块写入大文件内容 (Go默认使用UTF-8编码)
	f, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("创建缓存文件失败: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return "", "", fmt.Errorf("写入缓存文件失败: %w", err)
	}

	// 确保所有内容都写入磁盘
	err = f.Sync()
	if err != nil {
		return "", "", fmt.Errorf("同步缓存文件失败: %w", err)
	}

	return id, filePath, nil
}

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

// 获取文件扩展名集合
func getFileExtensions(paths []string) []string {
	extMap := make(map[string]bool)
	
	for _, path := range paths {
		ext := filepath.Ext(path)
		if ext != "" {
			extMap[ext] = true
		}
	}
	
	// 将map转换为切片
	var extensions []string
	for ext := range extMap {
		extensions = append(extensions, ext)
	}
	
	// 排序扩展名，使输出更易读
	sort.Strings(extensions)
	
	return extensions
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

	if len(filePaths) == 0 {
		fmt.Println("没有匹配的文件")
		os.Exit(0)
	}

	rootPath, formattedContent, relativePaths := processFiles(filePaths)

	fmt.Printf("Root path: %s\n\n", rootPath)
	fmt.Printf("%d file(s):\n", len(relativePaths))
	for _, path := range relativePaths {
		fmt.Printf("  - %s\n", path)
	}

	// 分别处理剪贴板和缓存
	contentSize := len(formattedContent)
	fmt.Printf("内容大小: %d 字节\n", contentSize)

	// 先保存到缓存文件
	id, filePath, err := saveToCache(formattedContent + "\n\n")
	if err != nil {
		fmt.Printf("保存到缓存文件失败: %v\n", err)
	} else {
		fmt.Printf("Cache ID: %s\n", id)
		fmt.Printf("Cache file: %s\n", filePath)
	}

	// 然后尝试写入剪贴板
	err = clipboard.WriteAll(formattedContent + "\n\n")
	if err != nil {
		fmt.Printf("复制到剪贴板失败: %v\n", err)
		fmt.Println("内容太大无法复制到剪贴板，但已保存到缓存文件")
	} else {
		fmt.Println("内容已复制到剪贴板")
	}
	
	// 输出文件扩展名集合
	extensions := getFileExtensions(relativePaths)
	fmt.Print("文件类型: (")
	for i, ext := range extensions {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(ext)
	}
	fmt.Println(")")
}
