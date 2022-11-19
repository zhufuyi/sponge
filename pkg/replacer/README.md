## replacer

一个替换目录下文件内容库，支持本地目录下文件和通过embed嵌入目录文件替换。

<br>

### 使用示例

```go
//go:embed dir
var fs embed.FS

func demo(){
	//r, err := replacer.New("dir")
	//if err != nil {
	//	panic(err)
	//}
	r, err := replacer.NewWithFS("dir", fs)
	if err != nil {
		panic(err)
	}

	ignoreFiles := []string{}
	fields := []replacer.Field{
		{
			Old: "1234",
			New: "8080",
		},
		{
			Old:             "abcde",
			New:             "hello",
			IsCaseSensitive: true,  // abcde-->hello, Abcde-->Hello
		},
	}
	r.SetSubDirs(subPaths...)          // 只处理指定子目录，优先级最高
	r.SetIgnoreDirs(ignoreDirs...)   // specify the directory in the subdirectory where processing is ignored
	r.SetIgnoreFiles(ignoreFiles...)   // specify the files in the subdirectory to be ignored for processing
	r.SetReplacementFields(fields)   // 设置替换文本
	r.SetOutPath("", "test")             // 设置输出目录，如果为空，根据名称和时间生成文件输出文件夹
	err = r.SaveFiles()                   // 保存替换后文件
	if err != nil {
		panic(err)
	}

	fmt.Printf("save files successfully, out = %s\n", replacer.GetOutPath())
}
```