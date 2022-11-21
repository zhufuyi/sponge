## replacer

A library of replacement file content, supports replacement of files in local directories and embedded directory files via embed.

<br>

### Example of use

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
    subDirs := []string{"testDir/replace"}
    subFiles := []string{"testDir/foo.txt"}
    ignoreDirs := []string{"testDir/ignore"}
    ignoreFiles := []string{"test.txt"}
    fields := []Field{
        {
            Old: "1234",
            New: "....",
        },
        {
            Old:             "abcdef",
            New:             "hello_",
            IsCaseSensitive: true,
        },
    }
    r.SetSubDirsAndFiles(subDirs, subFiles...) // process only specified subdirectories and files
	r.SetIgnoreDirs(ignoreDirs...)   // specify the directory in the subdirectory where processing is ignored
	r.SetIgnoreFiles(ignoreFiles...)   // specify the files in the subdirectory to be ignored for processing
	r.SetReplacementFields(fields)   // set replacement fields
	r.SetOutPath("", "test")             // set output directory, if empty, generate file output folder based on name and time
	err = r.SaveFiles()                   // save the replaced file
	if err != nil {
		panic(err)
	}

	fmt.Printf("save files successfully, out = %s\n", replacer.GetOutPath())
}
```
