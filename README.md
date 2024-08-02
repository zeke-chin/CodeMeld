# CodeMeld

CodeMeld 是一个简单的命令行工具，旨在帮助开发者将多个代码文件整合成一个易于大型语言模型（LLM）处理的格式。它自动提取共同的根路径，并将每个文件的内容格式化为 Markdown 代码块，附带相对路径信息。

❌ 路径不支持空格 ❌

## 为什么使用 CodeMeld？

在与大型语言模型（如 GPT-3、GPT-4 等）协作时，通常需要提供足够的上下文信息。CodeMeld 允许您快速整合多个相关的代码文件，使其易于复制粘贴到与 LLM 的对话中。这样可以：

1. 提供更完整的项目上下文
2. 保持文件结构的清晰性
3. 减少手动格式化的时间
4. 提高与 LLM 协作的效率

## 功能

- 自动检测多个文件的共同根路径
- 将文件内容转换为 Markdown 格式的代码块
- 为每个代码块添加相对路径信息
- 支持处理包含空格的文件路径

## 安装

mac / linux 直接可以使用 shell 命令安装 windows 未测试

```sh
wget https://raw.githubusercontent.com/zeke-chin/CodeMeld/main/install.sh -O ~/codemeld.install.sh && chmod +X ~/codemeld.install.sh && sh ~/codemeld.install.sh && rm ~/codemeld.install.sh && codemeld -v
```

国内镜像

```sh
wget https://mirror.ghproxy.com/https://raw.githubusercontent.com/zeke-chin/CodeMeld/main/install.sh -O ~/codemeld.install.sh && chmod +X ~/codemeld.install.sh && sh ~/codemeld.install.sh && rm ~/codemeld.install.sh && codemeld -v
```

## 使用方法
```shell
NAME:
   codemeld - 一个将code内容格式化成LLM更友好的格式的剪切板命令行工具

USAGE:
   codemeld [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --files value, -f value   文件、文件夹目录 可以使用空格或者换行隔开
   --reg value, -r value     匹配文件后缀 使用空格或者换行隔开
   --ireg value, --ir value  忽略的文件后缀 使用空格或者换行隔开
   --include-hidden, --ih    是否查看隐藏文件夹下的文件 (default: false)
   --help, -h                show help
```
下面是一些 `codemeld` 工具的使用实例：

1. **指定单个文件**：
   ```sh
   codemeld --files "myfile.txt"
   ```

2. **指定多个文件**：
   ```sh
   codemeld --files "file1.txt file2.txt"
   ```

3. **指定文件夹目录**：
   ```sh
   codemeld --files "/path/to/directory"
   ```

4. **指定文件夹目录和文件**：
   ```sh
   codemeld --files "/path/to/directory file3.txt"
   ```

5. **使用正则匹配文件后缀**：
   ```sh
   codemeld --files "/path/to/directory" --reg ".go .py"
   ```

6. **忽略特定文件后缀**：
   ```sh
   codemeld --files "/path/to/directory" --ireg ".log .tmp"
   ```

7. **查看隐藏文件夹中的文件**：
   ```sh
   codemeld --files "/path/to/directory" --include-hidden
   ```

8. **结合所有选项**：
   ```sh
   codemeld --files "/path/to/directory file1.txt" --reg ".go .py" --ireg ".log" --include-hidden
   ```


接下来你就可以直接将剪切板中的内容 粘贴给LLM

后面接你需要询问的问题即可

![image-20240802163120303](./assets/image-20240802163120303.png)
