
### **Major Feature Updates**

1. **Customizable Primary Keys for SQL Tables**:
    - Primary key names are no longer restricted to `id`, and the type is no longer limited to integers. Other names and string types are now supported.

2. **Enhanced Protobuf Support**:
    - Added a `protoc` plugin for converting Protobuf to JSON.

3. **Improved Code Generation**:
    - Added a command for generating code based on custom templates and fields.
    - Added a command for generating code based on custom templates and SQL.
    - Added a command for generating code based on custom templates and Protobuf.
    - Introduced a web interface for generating code using custom templates.

### **Framework and Code Enhancements**

4. **Simplified service code generation**:
    - Removed default code blocks for service registration and discovery and Nacos configuration center. If needed, users can add them manually.

5. **Directory Structure Optimization**:
    - Moved `internal/model/init.go` to the `internal/database` directory.

6. **Simplified Dependencies**:
    - Replaced `pkg/ggorm` with `pkg/sgorm`, reducing code dependencies during compilation.
    - Delete dropped library `pkg/mysql`

### **New Commands and Tools**

7. **Command Enhancements**:
    - The `make proto` command now automatically initializes the database and imports dependencies from `types.proto`.

8. **Architecture Diagram Generation**:
    - Automatically generates project business architecture diagrams based on service configurations. e.g. https://github.com/zhufuyi/spograph/blob/main/example.png

### **Bug Fixes**

- [#78](https://github.com/zhufuyi/sponge/issues/78)
- [#83](https://github.com/zhufuyi/sponge/issues/83)
- [#86](https://github.com/zhufuyi/sponge/issues/86)
- [#88](https://github.com/zhufuyi/sponge/issues/88)
