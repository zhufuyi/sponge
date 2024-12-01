
### **Major Feature Updates**

1. **Customizable Primary Keys for SQL Tables**:
    - Primary key names are no longer restricted to `id`, and the type is no longer limited to integers. Other names and string types are now supported.

2. **Improved Code Generation**:
    - Added a command for generating code based on custom templates and fields.
    - Added a command for generating code based on custom templates and SQL.
    - Added a command for generating code based on custom templates and Protobuf.
    - Added a `protoc` plugin for converting Protobuf to JSON.
    - Introduced a web interface for generating code using custom templates.

### **Framework and Code Enhancements**

3. **Simplified service code generation**:
    - Removed default code blocks for service registration and discovery and Nacos configuration center. If needed, users can add them manually.

4. **Directory Structure Optimization**:
    - Moved `internal/model/init.go` to the `internal/database` directory.

5. **Simplified Dependencies**:
    - Replaced `pkg/ggorm` with `pkg/sgorm`, reducing code dependencies during compilation.
    - Delete dropped library `pkg/mysql`

### **New Commands and Tools**

6. **Command Enhancements**:
    - The `make proto` command now automatically initializes the database and imports dependencies from `types.proto`.

7. **Simplified command**:
    - Merge `sponge configmap` into `sponge config` and rename it `cm`, see the help `sponge config cm -h`.

8. **Architecture Diagram Generation**:
    - Automatically generates project business architecture diagrams based on service configurations. e.g. https://github.com/zhufuyi/spograph/blob/main/example.png

### **Upgrade Dependency Library Version**

- google.golang.org/grpc: `v1.61.0` --> `v1.67.1`
- github.com/grpc-ecosystem/go-grpc-middleware: `v1.3.0` --> `v2.1.0`
- github.com/redis/go-redis/v9: `v9.6.1` --> `v9.7.0`
- go.etcd.io/etcd/client/v3: `v3.5.4` --> `v3.5.13`

### **Bug Fixes**

- [#78](https://github.com/zhufuyi/sponge/issues/78)
- [#83](https://github.com/zhufuyi/sponge/issues/83)
- [#86](https://github.com/zhufuyi/sponge/issues/86)
- [#88](https://github.com/zhufuyi/sponge/issues/88)
