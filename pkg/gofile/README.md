## gofile

File and directory management libraries.

<br>

## Example of use

```go
    // determine if a file or folder exists
    gofile.IsExists("/tmp/test/")

    // get the path to program execution
    gofile.GetRunPath()

    // get all files in a directory (absolute path)
    gofile.ListFiles("/tmp/")

    // get all files in a directory by prefix (absolute path)
    gofile.ListFiles(dir, WithPrefix("READ"))

    // get all files in a directory by suffix (absolute path)
    gofile.ListFiles(dir, WithSuffix(".go"))

    // get all files in a directory based on a string (absolute path)
    gofile.ListFiles(dir, WithContain("file"))
```
