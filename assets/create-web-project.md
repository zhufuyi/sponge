## 3 Quickly create web projects

### 3.1 Create http service from mysql

#### 3.1.1 Creating a table

To generate code based on mysql's data table, first prepare a mysql service ([docker install mysql](https://github.com/zhufuyi/sponge/blob/main/test/server/mysql/docker-compose.yaml)). For example, mysql has a database school and a data table teacher under the database, as shown in the following sql.

```sql
CREATE DATABASE IF NOT EXISTS school DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

use school;

create table teacher
(
    id bigint unsigned auto_increment
        primary key,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    name varchar(50) not null comment 'username',
    password varchar(100) not null comment 'password',
    email varchar(50) not null comment 'mail',
    phone varchar(30) not null comment 'mobile phone number',
    avatar varchar(200) null comment 'avatar',
    gender tinyint not null comment 'gender, 1:male, 2:female, other values:unknown',
    age tinyint not null comment 'age',
    birthday varchar(30) not null comment 'date of birth',
    school_name varchar(50) not null comment 'school name',
    college varchar(50) not null comment 'college',
    title varchar(10) not null comment 'title',
    profile text not null comment 'personal profile'
)
    comment 'teacher';

create index teacher_deleted_at_index
    on teacher (deleted_at);
```

Import the SQL DDL into mysql to create a database school, with a table teacher under school.

<br>

#### 3.1.2 Generating http service code

Open a terminal and execute the command.

```bash
sponge web http \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --out=./edusys
```

Check the parameter description command `sponge web http -h`, note that the parameter **repo-addr** is the image repository address, if you use the [official docker image repository](https://hub.docker.com/), you only need to fill in the username of the registered docker repository, if you use the private repository address, you need to fill in the full repository address.

<br>

Generating the complete http service code is in the current directory edusys with the following directory structure.

```
.
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── cache
│    ├── config
│    ├── dao
│    ├── ecode
│    ├── handler
│    ├── model
│    ├── routers
│    ├── server
│    └── types
└── scripts
```

The Makefile file in the edusys directory, which integrates commands related to compiling, testing, running, and deploying, switches to the edusys directory to run the service at

```bash
# Update swagger documentation
make docs

# Compile and run services
make run
```

Copy http://localhost:8080/swagger/index.html to your browser to test the CRUD interface, as shown in Figure 3-1.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-swag.jpg)

* Figure 3-1 http swagger documentation interface*

<br>

By default, the service is only enabled for the metrics collection interface, per-minute resource statistics information, and other service governance is off by default. In practical applications, some adjustments are made as needed.

- To use redis as a cache, open the configuration file `configs/edusys.yml`, change the **cacheType** field value to redis, and fill in the **redis** configuration address and port.
- By default, the flow limiting, fusion, link tracking, service registration and discovery functions are off, you can open the configuration file `configs/edusys.yml` to turn on the relevant functions, if you turn on the link tracking function, you must fill in the jaeger configuration information; if you turn on the service registration and discovery function, you must fill in one of the consul, etcd, nacos configuration information.
- If a configuration field name is added or modified, execute the command `sponge config --server-dir=./edusys` to update the corresponding go struct; it is not necessary to execute the update command to modify only the field values.
- Modify the error code information corresponding to the CRUD interface, open `ingernal/ecode/teacher_http.go`, modify the variable **teacherNO** value, which is the only value that does not repeat, the return message description is modified according to your needs, the interface custom error codes for the teacher table operations are added here.

<br>

#### 3.1.3 Generating Handler Code

In a service, there is usually more than one data table, so if a new data table is added, how can the generated handler code be automatically populated into the existing service code, using the `sponge web handler` command, for example if two new data tables **course** and **teach** are added.

```sql
create table course
(
    id bigint unsigned auto_increment
        primary key,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    code varchar(10) not null comment 'course code',
    name varchar(50) not null comment 'course name',
    credit tinyint not null comment 'credits',
    college varchar(50) not null comment 'college',
    semester varchar(20) not null comment 'semester',
    time varchar(30) not null comment 'class time',
    place varchar(30) not null comment 'place of class'
)
    comment 'course';

create index course_deleted_at_index
    on course (deleted_at);


create table teach
(
    id bigint unsigned auto_increment
        primary key,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    teacher_id bigint not null comment 'teacher id',
    teacher_name varchar(50) not null comment 'teacher name',
    course_id bigint not null comment 'course id',
    course_name varchar(50) not null comment 'course name',
    score char(5) not null comment 'students evaluate the quality of teaching, 5 grades: A,B,C,D,E'
)
    comment 'teacher course';

create index teach_course_id_index
    on teach (course_id);

create index teach_deleted_at_index
    on teach (deleted_at);

create index teach_teacher_id_index
    on teach (teacher_id);
```

<br>

Generate handler code that contains CRUD business logic.

```bash
sponge web handler \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=course,teach \
  --out=./edusys
```

Check the parameter description command `sponge web handler -h`, the parameter `out` is to specify the existing service folder edusys, if the parameter `out` is empty, you must specify the `module-name` parameter, generate the handler submodule code in the current directory, then copy the handler code to the folder edusys, the effect of both ways are The effect is the same.

After executing the command, the course and teach related code is generated in the `edusys/internal` directory.

```
.
└── internal
      ├── cache
      ├── dao
      ├── ecode
      ├── handler
      ├── model
      ├── routers
      └── types
```

<br>

Switch to the edusys directory and execute the command to run the service.

```bash
# Update swagger documentation
make docs

# Compile and run services
make run
```

Copy http://localhost:8080/swagger/index.html to your browser to test the CRUD interface, as shown in Figure 3-2.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-swag2.jpg)

* Figure 3-2 http swagger documentation interface*

The actual use requires modifying the custom CRUD interface to return error codes and messages, opening the file `ingernal/ecode/course_http.go` to modify the variable **courseNO** value, and opening the file `ingernal/ecode/teach_http.go` to modify the variable **teachNO** values.

Although the CRUD interface of each data table is generated, it is not necessarily suitable for the actual business logic, so you need to add the business logic code manually, fill in the database operation code to the `internal/dao` directory, and the business logic code to the `internal/handler` directory.

<br>

### 3.2 Create http service from proto file

If the standard CRUD interface http service code is not required, you can customize the interface in the proto file and use the spong command to generate the http service and interface template code.

#### 3.2.1 Custom interfaces

The following is a sample file, teacher.proto, where each method defines the description of the route and swagger document. The tag and validate descriptions are added to the message as needed.

```protobuf
syntax = "proto3";

package api.edusys.v1;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option go_package = "edusys/api/edusys/v1;v1";

// Generate *.swagger.json basic information
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  host: "localhost:8080"
  base_path: ""
  info: {
    title: "edusys api docs";
    version: "v0.0.0";
  };
  schemes: HTTP;
  schemes: HTTPS;
  consumes: "application/json";
  produces: "application/json";
};

service teacher {
  rpc Register(RegisterRequest) returns (RegisterReply) {
    // Set up routing
    option (google.api.http) = {
      post: "/api/v1/Register"
      body: "*"
    };
    // Set the swagger document corresponding to the route
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "Registered Users",
      description: "Submit information for registration",
      tags: "teacher",
    };
  }

  rpc Login(LoginRequest) returns (LoginReply) {
    option (google.api.http) = {
      post: "/api/v1/login"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "Login",
      description: "Login",
      tags: "teacher",
    };
  }
}

message RegisterRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message RegisterReply {
  int64 id = 1;
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message LoginReply {
  string token = 1;
}
```

<br>

#### 3.2.2 Generating http service code

Open a terminal and execute the command.

```bash
sponge web http-pb \
  --module-name=edusys \
  --server-name=edusys \
  --project-name=edusys \
  --repo-addr=zhufuyi \
  --protobuf-file=./teacher.proto \
  --out=./edusys
```

Check the parameter description command `sponge web http-pb -h`, which supports \* sign matching (example `--protobuf-file=*.proto`), indicating that code is generated based on a bulk proto file, and multiple proto files include at least one service, otherwise code generation is not allowed.

The directory for generating http service code is shown below, there are some differences with the http service code directory generated by `sponge web http`, the new proto file related **api** and **third_party** directories are added, there is no **cache**, **dao**, **model**, **types** directories in the internal directory. **handler**, **types** directories, where **handler** is the directory that holds business logic template code, which will be automatically generated by command.

```
.
├── api
│    └── edusys
│          └──v1
├── build
├── cmd
│    └── edusys
│          └── initial
├── configs
├── deployments
│    ├── docker-compose
│    └── kubernetes
├── docs
├── internal
│    ├── config
│    ├── ecode
│    ├── routers
│    └── server
├── scripts
└── third_party
```

Switch to the edusys directory and execute the command to run the service.

```bash
# Generate *pb.go file, generate handler template code, update swagger documentation
make proto

# Compile and run services
make run
```

Copying http://localhost:8080/apis/swagger/index.html to the browser test interface, as shown in Figure 3-3, the request returns a 500 error because the template code (internal/handler/teacher_logic.go file) calls `panic("implement me")` directly, which is meant to prompt for business logic code to be filled in.

![sponge-framework](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/http-pb-swag.jpg)
*Figure 3-3 http swagger documentation interface*

<br>

#### 3.2.3 Adding a new interface

Depending on the business requirements, new interfaces need to be added, in two cases.

**(1) Add new interface to original proto file**

Open `api/edusys/v1/teacher.proto`, e.g. add the **bindPhone** method and fill in the routing and swagger documentation description information to finish adding a new interface.

Execution order.

```bash
# Generate *pb.go file, generate handler template code, update swagger documentation
make proto
```

Generate new template files in the `internal/handler` and `internal/ecode` directories, then copy the latest generated template code into the business logic code area at.

- The template code file with the suffix **.gen.datetime** is generated in the `internal/handler` directory (example teacher_logic.go.gen.xxxx225619), because teacher_logic.go already exists and will not overwrite the business logic code written, so a new file is generated. Open the file teacher_logic.go.gen.xxxx225619, copy the template code for the add method **bindPhone** interface to the teacher_logic.go file, and fill in the business logic code.
- The template code file with the suffix **.gen.datetime** is generated in the `internal/ecode` directory, and the **bindPhone** interface error code is copied into the teacher_http.go file.
- Delete all files with the suffix **.gen.datetime**.

<br>

**(2) Adding interfaces to new proto files**

For example, if a new **course.proto** file is added, the interface under **course.proto** must include the routing and swagger documentation description information, check **Chapter 3.2.1** and copy the **course.proto** file to the `api/edusys/v1` directory to complete the newly added interface.

Execution order.

```bash
# Generate *pb.go file, generate handler template code, update swagger documentation
make proto
```

Generate code files with the **course** name prefix in the `internal/handler`, `internal/ecode`, and `internal/routers` directories by doing the following two operations.

- Fill in the business code in the `internal/handler/course.go` file.
- Modify the custom error code and message description in the `internal/ecode/course_http.go` file.

<br>

#### 3.2.4 Refining the http service

The http service code generated by the `sponge web http-pb` command does not have code related to `dao`, `cache`, `model` and other manipulation data, users can implement it themselves, if you use mysql database and redis cache, you can use **sponge** tool to generate `dao`, `cache`, `model` code directly.

Generate CRUD operation database code command.

```bash
sponge web dao \
  --db-dsn=root:123456@(192.168.3.37:3306)/school \
  --db-table=teacher \
  --include-init-db=true \
  --out=./edusys
```

Check the parameter description command `sponge web dao -h`, the parameter `-include-init-db` is used only once in a service, remove the parameter `-include-init-db` the next time you generate `dao` code, otherwise it will result in not generating the latest `dao` code, because the db initialization code already exists.

Whether you implement the `dao` code yourself or use the `dao` code generated by sponge, there are a number of operations that need to be done afterwards.

- Add mysql and redis to the initialization and release resource code of the service, open the `cmd/edusys/initial/initApp.go` file, backcomment out the call to mysql and redis initialization code, open the `cmd/edusys/initial/registerClose.go` file , backcomment out the call to mysql and redis release resource code, the initial code is a one-time change.
- The generated `dao` code, and custom methods **register** and **login** can not correspond exactly, you need to manually in the file `internal/dao/teacher.go` to supplement the code (file name teacher is the name of the table), and then in the `internal/handler/teacher. go` to fill in the business logic code (filename teacher is the name of the proto file), the business code returns the error using the error code defined in the `internal/ecode` directory, if the error message is returned directly, the requesting side will receive an UNKNOWN error message, that is, an undefined error message.
- The default uses local memory for caching, change it to use redis as cache, change the field **cacheType** value to redis in the configuration file `configs/edusys.yml`, and fill in the redis address and port.

Switching to the edusys directory to run the service again.

```bash
# Compile and run services
make run
```

Open http://localhost:8080/apis/swagger/index.html to request the interface again and it returns data properly.

<br>

### 3.3 Summary

There are two ways to generate http service code, mysql and proto files.

- According to mysql generated http service code includes CRUD interface code for each data table, subsequent addition of new interfaces, you can refer to the CRUD way to add business logic code, the newly added interfaces need to manually fill in the swagger description information.
- The http service generated from the proto file does not include the manipulation database code though, nor the CRUD interface logic code, which can be generated using the `sponge web dao` command to manipulate the database code as needed. With the addition of the new interface, in addition to generating handler template code, swagger documentation, route registration code, and error codes for the interface are automatically generated.

Both ways can complete the same http service interface, according to the actual application choose one of them, if you do backend management services, use mysql to produce CRUD interface code directly, you can write less code. For most need to customize the interface service, use the proto file way to generate the http service, this way is also more freedom, after writing the proto file, in addition to the business logic code, other code is generated through the plug-in.

<br><br>
