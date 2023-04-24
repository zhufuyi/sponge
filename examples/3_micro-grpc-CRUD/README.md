[**micro-grpc-CRUD 中文示例**](https://www.bilibili.com/read/cv23064432)

<br>

micro-grpc-CRUD source code generated using sponge, [**micro-grpc-CRUD.zip**](https://github.com/zhufuyi/sponge/tree/main/examples/3_micro-grpc-CRUD/micro-grpc-CRUD.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

### Quickly create a microservice project

Enter the Sponge UI interface, click on the left menu bar 【sql】--> 【RPC type】-->【Create rpc project】, fill in some parameters to generate a complete microservice project code.

The microservice code is mainly composed of commonly used libraries such as [grpc](https://github.com/grpc/grpc-go), [gorm](https://github.com/go-gorm/gorm), [go-redis](https://github.com/go-redis/redis), and also includes rpc client CRUD test code, common service governance function code, build deployment scripts, etc.

![micro-rpc](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc.png)

Switch to the user directory and run the command:

```bash
# Generate pb.go code
make proto

# Compile and start rpc service
make run
```

Use goland IDE to open user service code, enter the internal/service directory, open the `teacher_client_test.go` file, you can test CRUD methods here, similar to testing CRUD interfaces in swagger interface. Fill in parameters before testing and click the green button to test.

![micro-rpc-test](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/micro-rpc-test.png)

<br>

### Batch add CRUD code to rpc service

Enter the Sponge UI interface, click on the left menu bar 【sql】--> 【RPC type】-->【Generate service CRUD code】, select any number of tables to generate code, then move the generated CRUD code to the rpc service directory to complete batch addition of CURD interfaces in microservices without changing any code.

![micro-rpc-service](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-service.png)

Switch to user service directory and run command:

```bash
# Generate pb.go code
make proto

# Compile and start user service
make run
```

Use goland IDE, enter internal/service directory, open `teach_client_test.go` and `course_client_test.go` files to test CRUD methods.
