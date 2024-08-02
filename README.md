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

```sh
wget https://mirror.ghproxy.com/https://raw.githubusercontent.com/zeke-chin/CodeMeld/main/install.sh -O ~/codemeld.install.sh && chmod +X ~/codemeld.install.sh && sh ~/codemeld.install.sh && rm ~/codemeld.install.sh && codemeld -v
```

## 使用方法

在命令行中运行以下命令：

```
codemeld "文件路径1 文件路径2 文件路径3 ..."
```



例如：

```
codemeld "/path/to/file1.py /path/to/file2.js /path/to/file3.rs"

或者

codemeld "/path/to/file1.py
/path/to/file2.js
/path/to/file3.rs"
```

接下来你就可以直接粘贴在与LLM的文本输入框内

后面接你需要询问的问题即可

![image-20240802163120303](./assets/image-20240802163120303.png)

