[**micro-grpc-protobuf 中文示例**](https://www.bilibili.com/read/cv23099236)

<br>

micro-grpc-protobuf source code generated using sponge, [**micro-grpc-protobuf.zip**](https://github.com/zhufuyi/sponge/tree/main/examples/4_micro-grpc-protobuf/micro-grpc-protobuf.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

### Quickly create a microservice project

Prepare a proto file before creating a microservice, Enter the sponge UI interface, click on 【protobuf】--> 【RPC type】-->【Create RPC project】in the left menu bar, fill in some parameters to generate common microservice project code.

The microservice framework uses [grpc](https://github.com/grpc/grpc-go), and also includes commonly used service governance function code, build deployment scripts, etc. The database used is chosen by yourself.

![micro-rpc-pb](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-pb.png)

Switch to the user directory and run the command:

```bash
# Generate pb.go code, generate template code, generate test code
make proto

# Open internal/service/user.go, this is the generated template code. 
# There is a panic code that prompts you to fill in business logic code. Fill in business logic here.

# Compile and start user service
make run
```

Open user service code with goland IDE, go to internal/service directory, open `user_client_test.go` file, you can test rpc method here, similar to testing interface on swagger interface. Fill in parameters before testing and click green button to test.

![micro-rpc-pb-test](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/micro-rpc-pb-test.png)

<br>

### Batch add any api interface code embedded into microservices

Open `api/user/v1/user.proto` file and add 2 rpc methods for changing password and logging out. You can also add rpc methods in newly created proto files.

Switch to user service directory and run command:

```bash
# Generate pb.go code, generate template code, generate test code
make proto

# Go to internal/service/ directory, open file with date suffix, copy newly added interface code 
# to user.go file, remove panic code prompt code and fill in business logic

# Clear files with date suffix
make clean

# Compile and start user service
make run
```

Use goland IDE, go to internal/service directory, open `user_client_test.go` file, you can test newly added rpc methods here.
