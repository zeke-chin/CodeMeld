import os
import sys
import shlex
from pathlib import Path
import pyperclip

def get_common_root(paths):
    return str(Path(os.path.commonpath(paths))) if paths else None

def read_file_content(file_path):
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            return file.read()
    except Exception as e:
        return f"Error reading file: {str(e)}"

def process_files(file_paths):
    root_path = get_common_root(file_paths)
    result = []
    relative_paths = []
    for path in file_paths:
        content = read_file_content(path)
        relative_path = os.path.relpath(path, root_path)
        relative_paths.append(relative_path)
        result.append(f"```{relative_path}\n{content}\n```")
    return root_path, "\n\n".join(result), relative_paths

def main(input_string):
    try:
        file_paths = shlex.split(input_string)
        if not file_paths:
            raise ValueError("No file paths provided. / 未提供文件路径。")
        root_path, formatted_content, relative_paths = process_files(file_paths)
        print(f"Root path: {root_path}\n")
        print(f"{len(relative_paths)} File(s):")
        for path in relative_paths:
            print(f"  - {path}")
        print("\nFormatted content copied to clipboard.")
        return formatted_content
    except Exception as e:
        print(f"Error / 错误: {str(e)}")
        return None

def print_usage():
    print("用法 | Usage:\n python script.py \"path1 path2 path3 ...\"")
    print("\n描述 | Description:\n")
    print("  将多个文件的内容以LLM更加容易理解的方式格式化并复制到剪贴板。")
    print("  Format and copy the contents of multiple files to the clipboard in a way that is easier for LLMs to understand.")
    print("\n示例 | Example (可以直接使用GUI复制出的格式之间黏贴 | You can directly copy and paste between formats using the GUI.):\n")
    print("  python script.py \"./file1.txt /home/user/file2.py ~/file3.md\"")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print_usage()
        sys.exit(1)
    content = main(sys.argv[1])
    if content:
        pyperclip.copy(content)