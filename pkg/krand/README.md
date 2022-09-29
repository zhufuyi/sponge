## krand

生成随机字符串、整数、浮点数。

<br>

## 使用示例

### 生成随机字符串

```go
    /*
	R_NUM = 1      // R_NUM 纯数字
	R_UPPER = 2   // R_UPPER 大写字母
	R_LOWER = 4  // R_LOWER 小写字母
	R_All = 7	       // R_All 数字、大小写字母
    */

	// 通过|或组合出不同类型
    kind := krand.R_NUM|krand.R_UPPER    // 大写字母和数字
    // kind := krand.R_All    // 大小写字母和数字

    krand.String(kind, 10)   // 长度为10，大写字母和数字组成的随机字符串
```

<br>

### 生成随机整数

```go
    krand.Int(200)            // 随机数范围0 ~ 200
    krand.Int(1000, 2000)  // 随机数范围1000 ~ 2000
```

<br>

### 生成随机浮点数

```go
    krand.Float64(1, 200)            // 1位小数点的浮点数，范围0~200
    krand.Float64(2, 100,1000)            // 2位小数点的浮点数，范围100~1000
```
