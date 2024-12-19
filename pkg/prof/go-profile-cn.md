
### 通过http api采集profile

> [!note] 要能够使用http api采集profile，需要在`configs`目录下yaml文件设置 `enableMetrics: true`，默认路由是`/debug/pprof`。

在服务中通过http api采集profile[代码示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/prof/README.md#sampling-profile-by-http)。

通常在开发或测试时使用，如果线上开启会有一点点性能损耗，根据实际情况是否开启使用。除了支持go语言本身提供默认的profile分析，还支持io分析，路由是`/debug/pprof/profile-io`。

- 对于web服务，默认采集profile地址 http://localhost:8080/debug/pprof
- 对于grpc服务，默认采集profile地址 http://localhost:8283/debug/pprof

结合**go tool pprof**工具，任意时刻都可以分析当前程序运行状况。

<br>

### 通过系统信号通知采集profile

> [!note] 默认已开启系统信号通知采集profile功能，不需要额外配置。

在服务中开启系统信号通知采集profile[代码示例](https://github.com/go-dev-frame/sponge/blob/main/pkg/prof/README.md#sampling-profile-by-system-notification-signal)。

使用`http接口`方式，程序后台一直定时记录profile相关信息等，绝大多数时间都不会去读取这些profile，可以改进一下，只有需要的时候再开始采集profile，采集完后自动关闭，sponge生成的服务支持监听系统信号来开启和停止采集profile，默认使用了 **SIGTRAP**(5) 系统信号(linux环境建议改为SIGUSR1，windows环境不支持SIGUSR1)，发送信号给服务：

```bash
# 通过名称查看服务pid(第二列)
ps aux | grep 服务名称

# 发送信号给服务
kill -trap pid值
# kill -usr1 pid值
```

服务收到系统信号通知后，开始采集profile并保存到`/tmp/服务名称_profile`目录，默认采集时长为60秒，60秒后自动停止采集profile，如果只想采集30秒，发送第一次信号开始采集，大概30秒后再发送第二次信号停止采集profile，类似开关。默认采集**cpu**、**memory**、**goroutine**、**block**、**mutex**、**threadcreate**六种类型profile，文件格式`日期时间_pid_服务名称_profile类型.out`，示例：

```
xxx221809_58546_user_cpu.out
xxx221809_58546_user_mem.out
xxx221809_58546_user_goroutine.out
xxx221809_58546_user_block.out
xxx221809_58546_user_mutex.out
xxx221809_58546_user_threadcreate.out
```

因为trace的profile文件相对比较大，因此默认没有采集，根据实际需要可以开启采集trace(在服务初始化时调用`prof.EnableTrace()`)。

获得离线文件后，使用pprof工具使用交互式或界面方式进行分析：

```bash
# 交互式
go tool pprof [options] source

# 界面
go tool pprof -http=[host]:[port] [options] source
```

<br>

### 自适应采集profile

> [!note] 要能够使用系统信号来通知采集profile，需要在configs目录下yaml文件设置 `enableStat: true`

在线上运行的服务，没有出问题时基本不会去手动采集profile，但是又想在服务发出告警同时采集profile文件。为了解决这个问题，sponge创建的web或grpc务默认支持自适应采集profile功能，是把`系统信号通知采集profile`与`资源统计的告警功能`结合起来实现的，告警条件：

- 记录程序的cpu使用率连续3次(默认每分钟一次)，3次平均使用率超过80%时触发告警。
- 记录程序的物理内存使用率3次(默认每分钟一次)，3次平均占用系统内存超过80%时触发告警。
- 如果持续超过告警阈值，默认间隔15分钟发出一次告警。

触发告警时，程序内部调用kill函数发送系统信号通知采集profile，采集的profile文件保存到`/tmp/服务名称_profile`目录，即使在半夜程序的cpu或内存过高，第二天也可以通过分析profile来发现程序哪里造成cpu或内存过高。

> [!note] 自适应采集profile功能不支持windows环境。

<br>
