[**micro-gin-rpc-gateway 中文示例**](https://www.bilibili.com/read/cv23189890)

<br>

micro-gin-rpc-gateway source code generated using sponge, [**micro-gin-rpc-gateway.zip**](https://github.com/zhufuyi/sponge/tree/main/examples/5_micro-gin-rpc-gateway/micro-gin-rpc-gateway.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

### Quickly create an RPC gateway project

Before creating an RPC gateway project, prepare a .proto file. The .proto file must contain **routing description information** and **Swagger description information**.

Enter the Sponge UI interface, click the left menu bar 【protobuf】 → 【Web type】 → 【Create RPC gateway project】, fill in some parameters to generate the RPC gateway project code.

The web framework uses [gin](https://github.com/gin-gonic/gin), which also includes Swagger documentation, common service governance code, build deployment scripts, etc.

![micro-rpc-gw-pb](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-gw-pb.png)

In order to connect the RPC service in the RPC gateway service, you need to generate code to connect the RPC service separately. Click the left menu bar 【Public】 → 【Generate code to connect RPC services】, fill in some parameters to generate the code, and then move the generated code to connect RPC services to the RPC gateway project code.

![micro-rpc-cli](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-cli.png)

In the RPC gateway service, in order to call the methods of the RPC service, you need to copy the .proto file of the RPC service to the api/usergw/v1 directory of the RPC gateway service.

Switch to the usergw directory and execute the command:

```bash
# Generate pb.go code, generate route registration code, generate template code, and generate Swagger documentation 
make proto  

# Open internal/service/usergw_logic.go, this is the generated API interface code, 
# with a line of panic code prompting you to fill in business logic code. 
# First import the generated code to connect the RPC server, and call the RPC service method in the business logic.

# Compile and start the web service
make run
```
Open [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) in your browser to test the API interface.

![micro-rpc-gw-pb-swagger](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/micro-rpc-gw-pb-swagger.png)

<br>

If you need to add new API interfaces later, add RPC methods and messages to the proto file api/usergw/v1/usergw.proto.

If you need to connect other RPC services later, repeat the above steps to "generate code to connect RPC services". 
