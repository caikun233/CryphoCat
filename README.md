# CryphoCat
~~***该项目仍在alpha测试阶段***~~

~~***该项目仍在alpha测试阶段***~~

~~***该项目仍在alpha测试阶段***~~

***最近在玩Golang，所以用Go简单写了个命令行版本的**

简体中文ONLY

![](https://img.shields.io/badge/Golang-1.22.0-007d9c)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

## 介绍

**Go版就不写英文README了**

* 支持任意长度文本的非对称加密，目前正在使用RSA2048 （我想使其成为可选项，但目前开发精力不在这，先这样放着吧）。
* 已经内建了“两人聊天”场景下的密钥目录生成规则，你可以发起Pull Request来添加更多功能。
* 全部数据均离线处理并开源，不上传任何数据，与聊天软件完全无关。
* 只要你使用的聊天软件能够确保发送和接收的消息是相同的，那么你们的交流就永远不会被监听内容。
* **这辈子都不会支持音视频的，除非谁有办法把音视频变成文本。**
* 你可以通过将图片进行base64编码，再进行加密的方式来传送图片。开发计划中提到了加入base64编码图片的功能。 ~~啥b才这么发图片，你就不能用几句话描述一下你想表达什么，然后不用那些b图吗？~~

## 下载 & 使用

### 从源码运行

#### Linux

1. 安装Go环境（太多教程了不再赘述）
1. 下载源码文件（用Linux这还用教）
1. 从命令行运行

​		go run ./cryphocat4go.go

#### Windows

1. 同样安装Go环境，Golang官方有[安装包下载地址](https://golang.google.cn/dl/go1.22.3.windows-amd64.msi)，无脑下一步就行
2. 下载源码文件
3. 从命令行运行

​		go run cryphocat4go.go

### 直接运行可执行文件

在release里下载适合你系统的文件，Windows双击打开、Linux在终端上权限执行即可。

### 自行编译

#### Linux

1. 安装Go环境

2. 下载源码文件（用Linux这还用教）

3. 开始编译

   ```shell
   go build -o [你想输出的文件名] ./cryphocat4go.go
   ```

#### Windows

1. 同样安装Go环境，Golang官方有[安装包下载地址](https://golang.google.cn/dl/go1.22.3.windows-amd64.msi)，无脑下一步就行

2. 下载源码文件

3. 在.go源码文件的同级目录打开个终端，cmd或powershell都行

4. 开始编译

   ```shell
   go build .\cryphocat4go.go
   ```

5. 编译好的exe文件会与源码文件名一样，想在编译时修改的话，网上教程多的是

## 开发计划

- [ ] 使RSA密钥长度可选。
- [ ] 尝试加入对图片的base64编码。
- [ ] 受推友指点，可以添加一个计算图片哈希的功能，防止图片被平台添加追踪水印。
