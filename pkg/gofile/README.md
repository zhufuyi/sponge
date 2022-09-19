## gofile

对文件和目录管理库。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gofile

<br>

## 使用示例

```go
    // 判断文件或文件夹是否存在
    gofile.IsExists("/tmp/test/")

    // 获取程序执行的路径
    gofile.GetRunPath()

    // 获取目录下的所有文件(绝对路径)
    gofile.ListFiles("/tmp/")

    // 根据前缀获取目录下的所有文件(绝对路径)
    gofile.ListFiles(dir, WithPrefix("READ"))

    // 根据后缀获取目录下的所有文件(绝对路径)
    gofile.ListFiles(dir, WithSuffix(".go"))

    // 根据字符串获取目录下的所有文件(绝对路径)
    gofile.ListFiles(dir, WithContain("file"))
```
