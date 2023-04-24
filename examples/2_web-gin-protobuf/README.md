[web-gin-protobuf 中文示例](https://www.bilibili.com/read/cv23040234)

<br>

web-gin-protobuf source code generated using sponge, [web-gin-protobuf.zip](https://github.com/zhufuyi/sponge/tree/main/examples/2_web-gin-protobuf/web-gin-protobuf.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

### Quickly Create a Web Project

Prepare a proto file before creating a web service. The proto file must contain **route description information** and **swagger description information**.

Enter the UI interface of sponsor, click 【protobuf】-->【Web Type】-->【Create Web Project】in the left menu bar, and fill in some parameters to generate the web service project code.

The web framework uses [gin](https://github.com/gin-gonic/gin). It also includes swagger documents, common service governance function codes, and build and deployment scripts. You can choose which database to use.

![web-http-pb](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http-pb.png)

Change to the web directory and execute the command:

```bash
# Generate API code, generate registered routing code, and generate swagger docs 
make proto

# Open internal/handler/user_logic.go, which is the generated API code. There is 
# a line of panic code prompting to fill in the business logic code. Fill in the business logic here.
# The registration route of the API, the go structure code of the input parameter and the returned result, 
# the swagger document, and the definition error code have all been generated. 
# Just fill in the business logic code. 

# compile and start web services
make run 
```

Open [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) in your browser to test the API.

![web-http-pb-swagger](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http-pb-swagger.png)

<br>

### Add arbitrary API codes to embed into web services in batches

Add login and logout APIs to the proto file under the web service directory `api/user/v1`. You can also add APIs to the newly created proto file.

Enter the Sponge UI interface, click on the left menu bar 【sql】--> 【Web type】-->【Generate handler CRUD code】, select any number of tables to generate code, then move the generated CRUD code to the web service directory to complete batch addition of CURD interfaces in the web service without changing any code.

Change to the web directory and execute the command:

```bash
# Generate API code, generate registered routing code, and generate swagger documents
make proto 

# Enter the internal/handler/directory, open the file with the date suffix, copy the newly added 
# API code to the user_logic.go file, remove the panic code prompt code, and fill in the business logic 

#Remove files with date suffix
make clean 

#compile and start web services
make run 
```

Open [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) in your browser to test the new API.

![web-http-pb-swagger2](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http-pb-swagger2.png)
