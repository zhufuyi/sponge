[**web-gin-CRUD 中文示例**](https://www.bilibili.com/read/cv23018269)

<br>

web-gin-CRUD source code generated using sponge, [**web-gin-CRUD.zip**](https://github.com/zhufuyi/sponge/tree/main/examples/1_web-gin-CRUD/web-gin-CRUD.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

### Quickly create a web project

Enter the Sponge UI interface, click on the left menu bar 【sql】--> 【Web type】-->【Create web project】, fill in some parameters to generate the complete project code for the web service.

The web service code is mainly composed of commonly used libraries such as [gin](https://github.com/gin-gonic/gin), [gorm](https://github.com/go-gorm/gorm), [go-redis](https://github.com/go-redis/redis), and also includes swagger documentation, test code, common service governance function code, build deployment scripts, etc.

![web-http](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http.png)

Switch to the web directory and run the command:

```bash
# Generate swagger documentation
make docs

# Compile and start the web service
make run
```

Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser to perform CRUD operations on the table.

![web-http-swagger](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http-swagger.png)

<br>

### Batch add CRUD code embedded in web service

Enter the Sponge UI interface, click on the left menu bar 【sql】--> 【Web type】-->【Generate handler CRUD code】, select any number of tables to generate code, then move the generated CRUD code to the web service directory to complete batch addition of CURD interfaces in the web service without changing any code.

![web-http-handler](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_web-http-handler.png)

Switch to the web directory and run the command:

```bash
# Generate swagger documentation
make docs

# Compile and start the web service
make run
```

Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser to see the newly added CRUD interfaces.
