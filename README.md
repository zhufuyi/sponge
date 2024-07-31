<p align="center">
<img width="500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.png">
</p>

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

</div>

## 当前版本为魔改版本，具体用法请看 [官方文档](https://github.com/zhufuyi/sponge) 

### 源码安装 使用说明
```go
    git clone https://github.com/ice-leng/sponge.git
    cd sponge/cmd/sponge
    go run ./main.go init
```

### 主要魔改功能有
- 基于数据库dsn 添加表前缀 
```html
    数据库dsn: root:@(127.0.0.1:3306)/hyperf;prefix=t_
```
- 去掉下载代码功能，替换为，命令行 在那个目录，代码就在这个目录下生成
```go
    mkdir xxx
    cd xxx
    sponge run 
    ... // web 操作 代码下载 
	ls -al
```