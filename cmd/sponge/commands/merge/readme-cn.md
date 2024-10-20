## 合并命令

合并命令用于将生成的代码自动合并到已有的模板文件中，无需担心影响已编写的业务逻辑代码。如果合并过程中出现意外，合并前的代码备份会保存在 `/tmp/sponge_merge_backup_code` 目录中，您可以从中恢复之前的代码状态。

当自动合并代码出错时(proto文件中service数量变化导致)，需要手动合并代码。

<br>

### 手动合并代码说明

在大多数情况下，一个 proto 文件通常定义一个 service。但在开发过程中，proto 文件中的 service 数量可能会有所变化，例如从一个 service 增加到多个，或者从多个 service 减少为一个。这种变化可能导致自动合并代码失败，此时就需要手动进行代码合并。

#### 增加 Service 时的手动合并

以下以 `greeter.proto` 为例，初始文件中只有一个 service，名为 `Foobar1`。文件内容如下：

```protobuf
syntax = "proto3";

package greeter;

option go_package = "greeter";

service Foobar1 {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

在只有一个 service 的情况下，自动合并代码通常能正常工作，无需手动干预。

假设现在需要在 `greeter.proto` 中增加一个名为 `Foobar2` 的 service，更新后的文件内容如下：

```protobuf
syntax = "proto3";

package greeter;

option go_package = "greeter";

service Foobar1 {
  rpc SayHello (SayHelloRequest) returns (SayHelloReply) {}
}

message SayHelloRequest {
  string name = 1;
}

message SayHelloReply {
  string message = 1;
}

service Foobar2 {
  rpc SayHello (SayHelloRequest) returns (SayHelloReply) {}
}
```

当 proto 文件中新增 service 后，自动合并代码会导致失败，此时需要手动合并代码。手动合并的步骤如下：

1. 根据错误提示，找到生成的代码文件（文件后缀格式为 `.go.gen<日期时间>`），打开文件。

2. 找到 service `Foobar2` 对应的 Go 代码块，将其复制到目标 Go 文件中（即与 `.go.gen<日期时间>`文件相同的前缀go文件）。

3. 复制后的 Go 文件需符合以下要求：
    - Go文件的 service 代码块数量与 proto 文件的 service 一样，并且顺序必须一致。
    - Go文件的 service 代码块必须有固定的分割标记：`// ---------- Do not delete or move this split line, this is the merge code marker ----------`。
    - 当 proto 文件中仅剩一个 service 时，Go文件的 service 代码块不需要分割标记。

通过手动合并代码，后续如果 proto 文件中 service 数量不变化，自动合并功能都可以正常工作。

<br>

#### 减少 Service 时的手动合并

继续以 `greeter.proto` 为例，此时文件中包含两个 service：`Foobar1` 和 `Foobar2`。在这种情况下，自动合并代码通常能正常运行，无需手动干预。

假设现在需要删除 `Foobar2`，更新后的 `greeter.proto` 文件内容如下：

```protobuf
syntax = "proto3";

package greeter;

option go_package = "greeter";

service Foobar1 {
  rpc SayHello (SayHelloRequest) returns (SayHelloReply) {}
}

message SayHelloRequest {
  string name = 1;
}

message SayHelloReply {
  string message = 1;
}
```

此时，自动合并代码会导致失败，需手动合并代码。手动合并的步骤如下：

1. 根据错误提示，打开生成的代码文件（同样为 `.go.gen<时间戳>` 格式）。

2. 找到与 `Foobar2` 对应的 Go 代码块，删除该代码块。

3. 手动调整后go文件需满足以下要求：
    - Go文件的 service 代码块数量与 proto 文件的 service 一样，并且顺序必须一致。
    - Go文件的 service 代码块必须有固定的分割标记：`// ---------- Do not delete or move this split line, this is the merge code marker ----------`。
    - 当 proto 文件中仅剩一个 service 时，Go文件的 service 代码块不需要分割标记。

手动调整完成后，后续如果 proto 文件中 service 数量不变化，自动合并功能都可以正常工作。

<br>

### 警告

如果在一个 proto 文件中包含多个service, 一旦生成代码之后，不要调整proto文件中的service的顺序，否则会导致自动合并代码混乱，此时只能手动调整合并后的代码。
