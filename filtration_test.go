package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestProcessFiles(t *testing.T) {
	// 创建一个临时目录作为测试环境
	tempDir, err := ioutil.TempDir("", "test_dir")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 在临时目录中创建测试文件结构
	files := []string{
		"file1.txt",
		"file2.go",
		"file3.tmp",
		"subdir/file4.go",
		"subdir/file5.txt",
		".hidden/file6.go",
	}

	for _, file := range files {
		path := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatalf("无法创建目录: %v", err)
		}
		err = ioutil.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("无法创建文件: %v", err)
		}
	}

	tests := []struct {
		name           string
		args           AppArgs
		expectedFiles  []string
		expectedErrNil bool
	}{
		{
			name: "包含所有文件",
			args: AppArgs{
				Files:         []string{tempDir},
				Reg:           []string{},
				IReg:          []string{},
				IncludeHidden: true,
			},
			expectedFiles:  []string{"file1.txt", "file2.go", "file3.tmp", "subdir/file4.go", "subdir/file5.txt", ".hidden/file6.go"},
			expectedErrNil: true,
		},
		{
			name: "只包含 .go 文件",
			args: AppArgs{
				Files:         []string{tempDir},
				Reg:           []string{".go"},
				IReg:          []string{},
				IncludeHidden: false,
			},
			expectedFiles:  []string{"file2.go", "subdir/file4.go"},
			expectedErrNil: true,
		},
		{
			name: "排除 .tmp 文件",
			args: AppArgs{
				Files:         []string{tempDir},
				Reg:           []string{},
				IReg:          []string{".tmp"},
				IncludeHidden: false,
			},
			expectedFiles:  []string{"file1.txt", "file2.go", "subdir/file4.go", "subdir/file5.txt"},
			expectedErrNil: true,
		},
		{
			name: "包含隐藏文件",
			args: AppArgs{
				Files:         []string{tempDir},
				Reg:           []string{".go"},
				IReg:          []string{},
				IncludeHidden: true,
			},
			expectedFiles:  []string{"file2.go", "subdir/file4.go", ".hidden/file6.go"},
			expectedErrNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filtrationFiles(tt.args)

			// 检查错误
			if (err == nil) != tt.expectedErrNil {
				t.Errorf("processFiles() error = %v, expectedErrNil %v", err, tt.expectedErrNil)
				return
			}

			// 对结果进行排序，以确保比较的一致性
			sort.Strings(result)
			expected := make([]string, len(tt.expectedFiles))
			for i, file := range tt.expectedFiles {
				expected[i] = filepath.Join(tempDir, file)
			}
			sort.Strings(expected)

			// 比较结果
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("processFiles() = %v, want %v", result, expected)
			}
		})
	}
}
