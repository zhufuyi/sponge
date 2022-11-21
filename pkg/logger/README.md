## logger

Log library encapsulated in [zap](https://github.com/uber-go/zap).

- Support for terminal printing and log saving.
- Support for automatic log file cutting.
- Support for json format and console log format output.
- Supports Debug, Info, Warn, Error, Panic, Fatal, also supports fmt.Printf-like log printing, Debugf, Infof, Warnf, Errorf, Panicf, Fatalf.

<br>

## Example of use

```go
    // (1) used directly, it will be initialised by default
    logger.Info("this is info")
    logger.Warn("this is warn", logger.String("foo","bar"), logger.Int("size",10), logger.Any("obj",obj))
    logger.Error("this is error", logger.Err(err), logger.String("foo","bar"))

    // (2) Initialize and then use
    logger.Init(
        logger.WithLevel("info"),     // set the data logging level, the default is debug
        logger.WithFormat("json"),  // set output format, default console
        logger.WithSave(true,         // set whether to save the log locally, default false
        //    logger.WithFileName("my.log"),      // file name, default is "out.log"
        //    logger.WithFileMaxSize(5),              // maximum file size (MB), default 10
       //     logger.WithFileMaxBackups(5),        // maximum number of old files, default 100
       //     logger.WithFileMaxAge(10),             // maximum number of days for old documents, default 30
       //     logger.WithFileIsCompression(true), // whether to compress and archive old files, default false
        )
    )
```
