# CryphoCat
~~***该项目仍在alpha测试阶段***~~

~~***该项目仍在alpha测试阶段***~~

~~***该项目仍在alpha测试阶段***~~

***该项目即将使用其它GUI框架重制**

可以用我加密你的任何文字聊天喵~

目前仅支持简体中文和英文。

![](https://img.shields.io/badge/go-1.25-blue)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

**English Readme please go to [README.md](https://github.com/caikun233/CryphoCat/blob/main/README.md)**

* 支持任意长度文本的非对称加密，RSA 密钥长度可选（CLI：1024/2048/4096，默认 4096；GUI：2048）。
* 已经内建了"两人聊天"场景下的密钥目录生成规则，你可以发起 Pull Request 来添加更多功能。
* 全部数据均离线处理并开源，不上传任何数据，与聊天软件完全无关。
* 只要你使用的聊天软件能够确保发送和接收的消息是相同的，那么你们的交流就永远不会被监听内容。
* **这辈子都不会支持音视频的，除非谁有办法把音视频变成文本。**
* 你可以通过将图片进行 base64 编码后加密的方式来传送图片。

## 安装 & 使用

### Linux / Windows / macOS 命令行（CLI）🔨

1. 从 **Actions** 页面下载对应平台的 CLI 二进制文件，或自行从源码编译（见下文）。

2. 直接运行：

   ```shell
   ./CryphoCat_CLI_linux_amd64
   ```

   中文界面请加 `--lang zhcn` 参数：

   ```shell
   ./CryphoCat_CLI_linux_amd64 --lang zhcn
   ```

3. 程序将自动在当前目录下新建名为 `RSAkeys` 的文件夹，并在其中创建 `my/` 和 `friend/` 子目录来存放你的密钥对和对方的公钥。

### Linux / Windows 图形化界面（GUI）🔨

1. 从 **Actions** 页面下载对应平台的 GUI 二进制文件。**仅支持 x64，不提供 x86 支持。**
2. 双击运行（Windows），或在 Linux 下执行 `./CryphoCat_GUI_linux_amd64`。
3. 中文界面请加 `--lang zhcn` 参数：

   ```shell
   ./CryphoCat_GUI_linux_amd64 --lang zhcn
   ```

### 从源码编译 🛠

需要 **Go 1.25+**。Linux 下还需安装：
```shell
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev \
  libxi-dev libgl1-mesa-dev libxxf86vm-dev pkg-config
```

```bash
cd go-source

# GUI 版本（默认）：
go build -o CryphoCat .

# 仅 CLI 版本（无 GUI，无 OpenGL 依赖）：
go build -tags cli -o CryphoCat_CLI .

# 运行测试：
go test ./...
```

## 开发计划

- [x] 在源代码和 release 中加入英文支持。
- [x] 完成 GUI 版本的开发。
- [x] 使 RSA 密钥长度可选（CLI 支持 1024/2048/4096）。
- [x] 自动复制加密后的文本。
- [x] 计算并比对文件哈希（MD5 / SHA1 / SHA256）。
- [x] 用 Go + Fyne 重写（单一可执行文件，无需 Python）。
- [ ] 让 GUI 看起来好看点。
- [ ] 尝试加入对图片的 base64 编码。
