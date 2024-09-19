## English | [简体中文](readme-cn.md)

## Merge Command

The merge command is used to automatically merge the generated code into existing template files without affecting the existing business logic code. If any issues occur during the merging process, a backup of the code before merging will be saved in the `/tmp/sponge_merge_backup_code` directory, allowing you to restore the previous state of your code.

Manual merging of code is required when automatic merging fails due to changes in the number of services in the proto file.

<br>

### Manual Code Merging Instructions

In most cases, a proto file typically defines a single service. However, during development, the number of services in a proto file might change, such as increasing from one service to multiple services or decreasing from multiple services to one. Such changes may cause automatic merging to fail, necessitating manual code merging.

#### Manual Merging When Adding a Service

Let's take `greeter.proto` as an example, where the initial file contains a single service named `Foobar1`. The content of the file is as follows:

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

When there is only one service, the automatic code merge usually works fine without manual intervention.

Suppose you need to add a new service named `Foobar2` to `greeter.proto`. The updated file content is as follows:

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

After adding a new service to the proto file, the automatic code merge may fail, requiring manual code merging. The steps for manual merging are as follows:

1. Based on the error message, locate the generated code file (with a suffix format of `.go.gen<timestamp>`) and open the file.

2. Find the Go code block corresponding to the `Foobar2` service and copy it into the target Go file (the one with the same prefix as the `.go.gen<timestamp>` file).

3. The Go file after copying must meet the following requirements:
    - The number of service code blocks in the Go file must match the number of services in the proto file, and their order must be consistent.
    - The service code blocks in the Go file must be separated by a fixed marker: `// ---------- Do not delete or move this split line, this is the merge code marker ----------`.
    - If the proto file contains only one service, the service code block in the Go file does not need a separator.

By manually merging the code, the automatic merging feature will work correctly if the number of services in the proto file remains unchanged in the future.

<br>

#### Manual Merging When Removing a Service

Continuing with the `greeter.proto` example, the file currently contains two services: `Foobar1` and `Foobar2`. In this case, automatic code merging usually works fine without manual intervention.

Suppose you need to remove `Foobar2`, and the updated `greeter.proto` file is as follows:

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

At this point, automatic code merging may fail, requiring manual code merging. The steps for manual merging are as follows:

1. Based on the error message, open the generated code file (also with a `.go.gen<timestamp>` suffix format).

2. Locate the Go code block corresponding to `Foobar2` and delete this code block.

3. After manual adjustment, the Go file must meet the following requirements:
    - The number of service code blocks in the Go file must match the number of services in the proto file, and their order must be consistent.
    - The service code blocks in the Go file must be separated by a fixed marker: `// ---------- Do not delete or move this split line, this is the merge code marker ----------`.
    - If the proto file contains only one service, the service code block in the Go file does not need a separator.

After manual adjustments are complete, the automatic merging feature will work correctly if the number of services in the proto file remains unchanged in the future.

<br>

### Warning

If multiple services are included in a proto file, once the code is generated, do not adjust the order of the services in the proto file, otherwise it will cause confusion in the automatic merging code, and manual adjustment of the code is required.
