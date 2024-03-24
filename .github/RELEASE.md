## Change log

- Code generation supports multiple types of code repositories, you can choose to use `monolithic application single repository (monolith)`, `microservice multi-repository (multi-repo)`, or `microservice single repository (mono-repo)` according to your project needs.
- Added automated testing scripts for code generation commands.
- Based on protobuf to generate web services, the generated template code and documentation must meet the following conditions:
  - rpc cannot be set as stream type.
  - rpc must set http related information (router and method).
- RPC stream based on protobuf supports generating corresponding template code and client testing code.
- The generated code based on protobuf supports some common special types, such as Empty, Any, Timestamp, etc.
- Fixed known bugs.
