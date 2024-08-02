#!/bin/bash

# 函数：检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 函数：尝试安装到指定目录
try_install() {
    local install_dir="$1"
    if [ ! -d "$install_dir" ]; then
        mkdir -p "$install_dir"
    fi
    if [ -w "$install_dir" ]; then
        mv "$TMP_FILE" "$install_dir/codemeld"
        chmod +x "$install_dir/codemeld"
        echo "CodeMeld 已成功安装到 $install_dir/codemeld"
        return 0
    fi
    return 1
}

# 检查必要的命令是否存在
for cmd in curl jq; do
    if ! command_exists $cmd; then
        echo "错误: 未找到 $cmd 命令。请安装后再运行此脚本。" 1>&2
        exit 1
    fi
done

# 获取系统类型和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 将架构映射到 Github release 使用的命名
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "不支持的架构: $ARCH" 1>&2
        exit 1
        ;;
esac

# 询问用户是否使用镜像源
read -p "是否使用镜像源(ghproxy)? (1.是[默认] 0.否[官方]): " use_mirror

# 如果用户输入为空或为1，则使用1作为默认值
if [ -z "$use_mirror" ] || [ "$use_mirror" -eq 1 ]; then
  use_mirror=1
else
  use_mirror=${use_mirror}
fi

# 获取最新 release 的下载 URL
LATEST_RELEASE_URL=$(curl -s https://api.github.com/repos/zeke-chin/CodeMeld/releases/latest | jq -r '.assets[] | select(.name | test("codemeld-'$OS'-'$ARCH'")) | .browser_download_url')
if [ -z "$LATEST_RELEASE_URL" ]; then
    echo "无法找到适合 $OS-$ARCH 的下载链接" 1>&2
    exit 1
fi

# 如果选择使用镜像，则修改下载链接
if [ "$use_mirror" = "1" ]; then
    LATEST_RELEASE_URL="https://mirror.ghproxy.com/$LATEST_RELEASE_URL"
    echo "使用镜像源下载: $LATEST_RELEASE_URL"
else
    echo "使用官方源下载: $LATEST_RELEASE_URL"
fi

# 下载文件
TMP_FILE=$(mktemp)
echo "正在下载 CodeMeld..."
if ! curl -L -o "$TMP_FILE" "$LATEST_RELEASE_URL"; then
    echo "下载失败" 1>&2
    rm -f "$TMP_FILE"
    exit 1
fi

# 尝试安装到不同的目录
echo "正在安装 CodeMeld..."
if try_install "/usr/local/bin" || try_install "$HOME/bin" || try_install "$HOME/.local/bin"; then
    # 安装成功
    # 检查 PATH 中是否包含安装目录
    if ! echo $PATH | grep -q "$HOME/bin" && ! echo $PATH | grep -q "$HOME/.local/bin" && ! echo $PATH | grep -q "/usr/local/bin"; then
        echo "警告: 安装目录可能不在你的 PATH 中。你可能需要添加它到你的 PATH。"
        echo "建议在你的 ~/.bashrc 或 ~/.zshrc 文件中添加以下行："
        echo "export PATH=\$PATH:$HOME/bin:$HOME/.local/bin:/usr/local/bin"
    fi
else
    echo "安装失败：无法找到可写的安装目录" 1>&2
    rm -f "$TMP_FILE"
    exit 1
fi

# 清理临时文件
rm -f "$TMP_FILE"