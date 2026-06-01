# CryphoCat

离线非对称加密工具，用于保护聊天消息的隐私。支持 RSA、ECC、Curve25519 和 Kyber 抗量子算法。单一可执行文件——双击启动原生 Windows GUI，命令行下自动切换为终端交互模式。

![](https://img.shields.io/badge/go-1.22-blue)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

[English README](README.md)

## 主要特性

- **10 种加密算法**：RSA-2048/3072/4096、ECC P-256/P-384/P-521、Curve25519、Kyber-512/768/1024
- **紧凑输出**：Base64 + zlib 压缩编码，密文长度比旧版十六进制格式缩减约 60%
- **内存优先**：密钥默认仅存储在内存中，可选择写入磁盘
- **图片支持**：可加载图片自动编码为 base64 再加密；解密后自动识别图片并弹出预览
- **剪贴板集成**：自己的公钥一键复制；对方的公钥可直接从剪贴板导入
- **文件哈希**：内建 MD5 / SHA1 / SHA256 计算与校验
- **Windows 原生 GUI**，体积约 7.7 MB
- **10 种界面语言**：EN、简体中文、繁體中文(臺灣/香港)、日本語、한국어、Русский、Français、Deutsch、Español
- 完全离线，无任何数据上传

## 快速开始

从 [Releases](https://github.com/caikun233/CryphoCat/releases) 下载。

```
cryphocat.exe                    # 自动检测：终端 → CLI，双击 → GUI
cryphocat.exe --gui              # 强制启动 GUI
cryphocat.exe --cli --lang zhcn  # 强制 CLI，中文界面
```

密钥默认存储在 `.cryphocat_keys/` 目录，勾选"仅存储在内存中"后不落盘。

## 编译

需要 Go 1.22+。先安装 rsrc：`go install github.com/akavel/rsrc@latest`。

```
cd go-source
rsrc -manifest cryphocat.manifest -ico favicon.ico -o rsrc.syso
go build -ldflags="-s -w" -o cryphocat.exe .
go test ./...                     # 运行测试
```

非 Windows 平台（仅 CLI）：跳过 rsrc，直接 `go build`。

## 发布流程

推送版本标签即可触发自动构建和 GitHub Release：

```
git tag v1.0.0
git push origin v1.0.0
```

## 开发计划

- [x] RSA、ECC、Curve25519、Kyber 多算法加解密
- [x] 文件哈希计算与校验（MD5/SHA1/SHA256）
- [x] 密文自动复制到剪贴板；密钥剪贴板导入导出
- [x] 10 语言 i18n，Windows 系统语言自动检测
- [x] Windows 原生 GUI（~7.7 MB，从 42 MB 缩减）
- [x] CLI / GUI 启动自动检测
- [x] 标签推送触发 GitHub Release，自动生成发布说明
- [x] 图片加载为 base64 后加密，解密后自动预览
- [x] Base64+zlib 紧凑输出编码
- [x] 内存密钥存储
