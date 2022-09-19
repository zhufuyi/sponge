## logger

在[zap](https://github.com/uber-go/zap)封装的日志库。

- 支持终端打印和保存日志。
- 支持日志文件自动切割。
- 支持json格式和console日志格式输出。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/logger

<br>

## 使用示例

支持Debug、Info、Warn、Error、Panic、Fatal，也支持类似fmt.Printf打印日志，Debugf、Infof、Warnf、Errorf、Panicf、Fatalf

```go
    // (1) 直接使用，默认会初始化
    logger.Info("this is info")
    logger.Warn("this is warn", logger.String("foo","bar"), logger.Int("size",10), logger.Any("obj",obj))
    logger.Error("this is error", logger.Err(err), logger.String("foo","bar"))

    // (2) 初始化后再使用
    logger.Init(
        logger.WithLevel("info"),     // 设置数据日志级别，默认是debug
        logger.WithFormat("json"),  // 设置输出格式，默认console
        logger.WithSave(true,         // 设置是否保存日志到本地，默认false
        //    logger.WithFileName("my.log"),      // 文件名称，默认"out.log"
        //    logger.WithFileMaxSize(5),              // 最大文件大小(MB)，默认10
       //     logger.WithFileMaxBackups(5),        // 旧文件的最大个数，默认100
       //     logger.WithFileMaxAge(10),             // 旧文件的最大天数，默认30
       //     logger.WithFileIsCompression(true), // 是否压缩归档旧文件，默认false
        )
    )
```
