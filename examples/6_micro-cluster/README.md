[**micro-cluster 中文示例**](https://www.bilibili.com/read/cv23255594)

<br>

micro-cluster source code generated using sponge, [**micro-cluster.zip**](https://github.com/zhufuyi/sponge/tree/main/examples/6_micro-cluster/micro-cluster.zip) code file in the current directory, it is generated according to the following steps.

<br>
<br>

By taking a simple e-commerce microservice cluster as an example, the product details page contains product information, inventory information, and product evaluation information. These data are scattered in different microservices. The RPC gateway service assembles the required data and returns it to the product details page, as shown in the following figure:

![micro-cluster](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-cluster.png)

Four proto files have been prepared in advance. Each proto file generates corresponding service code:

- The comment.proto file defines RPC methods to obtain comment data based on the product ID and is used to generate the comment RPC service code.
- The inventory.proto file defines RPC methods to obtain inventory data based on the product ID and is used to generate the inventory RPC service code.
- The product.proto file defines RPC methods to obtain product details based on the product ID and is used to generate the product RPC service code.
- The shopgw.proto file defines RPC methods to assemble the data required by the product details page based on the product ID and is used to generate the shop RPC gateway service code.

<br>

### Quickly generate and start comment, inventory, and product microservices

#### Generate comment, inventory, and product microservice code

Enter the Sponge UI interface, click the left menu bar 【Protobuf】--> 【RPC type】-->【Create RPC project】, fill in the respective parameters of the comment, inventory, and product, and generate the comment, inventory, and product service code respectively.

The microservice framework uses [grpc](https://github.com/grpc/grpc-go) and also contains common service governance function codes, build deployment scripts, etc. You can choose whatever database you want.

Quickly create a comment service as shown below:

![micro-rpc-pb-comment](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-pb-comment.png)

Quickly create an inventory service as shown below:

![micro-rpc-pb-inventory](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-pb-inventory.png)

Quickly create a product service as shown below:

![micro-rpc-pb-product](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-pb-product.png)

Open three terminals, one for each of the comment, inventory, and product services.

<br>

#### Start the comment service

Switch to the comment directory and perform the following steps:
(1) Generate pb.go code, generate template code, and generate test code

```shell
make proto
```

(2) Open internal/service/comment.go, which is the generated template code. There is a line of panic code prompting you to fill in the business logic code. Fill in the business logic here, for example, fill in the return value:

```go
	return &commentV1.ListByProductIDReply{
		Total:     11,
		ProductID: 1,
		CommentDetails: []*commentV1.CommentDetail{
			{
				Id:       1,
				Username: "Mr Zhang",
				Content:  "good",
			},
			{
				Id:       2,
				Username: "Mr Li",
				Content:  "good",
			},
			{
				Id:       3,
				Username: "Mr Wang",
				Content:  "not good",
			},
		},
	}, nil
```

(3) Open the configs/comment.yml configuration file, find grpc, and modify the port and httpPort values under it:

```yaml
grpc:   
  port: 18203              # listen port   
  httpPort: 18213        # profile and metrics ports  
```

(4) Compile and start the comment service

```shell
make run
```

<br>

#### Start the inventory service

Switch to the inventory directory and perform the same steps as for the comment:

(1) Generate pb.go code, generate template code, and generate test code

```shell
make proto  
```

(2) Open internal/service/inventory.go, which is the generated template code. There is a line of panic code prompting you to fill in the business logic code. Fill in the business logic here, for example, fill in the return value:

```go
	return &inventoryV1.GetByIDReply{
		InventoryDetail: &inventoryV1.InventoryDetail{
			Id:      1,
			Num:     999,
			SoldNum: 111,
		},
	}, nil
```

(3) Open the configs/inventory.yml configuration file, find grpc, and modify the port and httpPort values under it:

```yaml
grpc:    
  port: 28203              # listen port  
  httpPort: 28213        # profile and metrics ports   
```

(4) Compile and start the inventory service

```shell
make run
```

<br>

#### Start the product service

Switch to the product directory and perform the same steps as for the comment:

(1) Generate pb.go code, generate template code, and generate test code

```shell
make proto
```

(2) Open internal/service/product.go, which is the generated template code. There is a line of panic code prompting you to fill in the business logic code. Fill in the business logic here, for example, fill in the return value:

```go
	return &productV1.GetByIDReply{
		ProductDetail: &productV1.ProductDetail{
			Id:          1,
			Name:        "Data cable",
			Price:       10,
			Description: "Android type c data cable",
		},
		InventoryID: 1,
	}, nil
```

(3) Open the configs/product.yml configuration file, find grpc, and modify the port and httpPort values under it:

```yaml
grpc:    
  port: 38203              # listen port  
  httpPort: 38213        # profile and metrics ports
```

(4) Compile and start the product service

```shell
make run
```

After the comment, inventory, and product microservices have started successfully, you can now generate and start the gateway service.

<br>

### Quickly generate and start the RPC gateway service

Enter the Sponge UI interface, click the left menu bar 【Protobuf】--> 【Web type】-->【Create RPC gateway project】, fill in some parameters to generate the RPC gateway project code.

The web framework uses [gin](https://github.com/gin-gonic/gin), which also contains swagger documentation, common service governance function codes, build deployment scripts, etc.

![micro-rpc-gw-pb-shopgw](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-rpc-gw-pb-shopgw.png)

In order to connect to the comment, inventory, and product RPC services, you need to generate additional connection RPC service codes. Click the left menu bar 【Public】--> 【Generate connection RPC service code】, fill in the parameters to generate the code, and then move the generated connection RPC service code to the RPC gateway project code directory.

![micro-rpc-cli](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/en_micro-cluster-rpc-cli.png)

In order to call the methods of the RPC services in the RPC gateway service, you need to copy the proto files of the comment, inventory, and product RPC services to the api/shopgw/v1 directory of the RPC gateway service.

Switch to the shopgw directory and perform the following steps:
(1) Generate pb.go code, generate registered routing code, generate template code, and generate swagger documentation

```shell
make proto
```

(2) Open internal/service/shopgw_logic.go, which is the generated API interface code. Fill in the business logic code here, fill in the following simple business logic code:

```go
package service

import (
	"context"

	commentV1 "shopgw/api/comment/v1"
	inventoryV1 "shopgw/api/inventory/v1"
	productV1 "shopgw/api/product/v1"
	shopgwV1 "shopgw/api/shopgw/v1"
	"shopgw/internal/rpcclient"
)

var _ shopgwV1.ShopGwLogicer = (*shopGwClient)(nil)

type shopGwClient struct {
	productCli   productV1.ProductClient
	inventoryCli inventoryV1.InventoryClient
	commentCli   commentV1.CommentClient
}

// NewShopGwClient creating rpc clients
func NewShopGwClient() shopgwV1.ShopGwLogicer {
	return &shopGwClient{
		productCli:   productV1.NewProductClient(rpcclient.GetProductRPCConn()),
		inventoryCli: inventoryV1.NewInventoryClient(rpcclient.GetInventoryRPCConn()),
		commentCli:   commentV1.NewCommentClient(rpcclient.GetCommentRPCConn()),
	}
}

func (c *shopGwClient) GetDetailsByProductID(ctx context.Context, req *shopgwV1.GetDetailsByProductIDRequest) (*shopgwV1.GetDetailsByProductIDReply, error) {
	productRep, err := c.productCli.GetByID(ctx, &productV1.GetByIDRequest{
		Id: req.ProductID,
	})
	if err != nil {
		return nil, err
	}

	inventoryRep, err := c.inventoryCli.GetByID(ctx, &inventoryV1.GetByIDRequest{
		Id: productRep.InventoryID,
	})
	if err != nil {
		return nil, err
	}

	commentRep, err := c.commentCli.ListByProductID(ctx, &commentV1.ListByProductIDRequest{
		ProductID: req.ProductID,
	})
	if err != nil {
		return nil, err
	}

	return &shopgwV1.GetDetailsByProductIDReply{
		ProductDetail:   productRep.ProductDetail,
		InventoryDetail: inventoryRep.InventoryDetail,
		CommentDetails:  commentRep.CommentDetails,
	}, nil
}
```

(3) Open the configs/shopgw.yml configuration file, find grpcClient, and add the addresses of the comment, inventory, and product RPC services:

```yaml
grpcClient:
  - name: "comment"
    host: "127.0.0.1"
    port: 18282
  - name: "inventory"
    host: "127.0.0.1"
    port: 28282
  - name: "product"
    host: "127.0.0.1"
    port: 38282
```

(4) Compile and start the shopgw service

```shell
make run
```

You can test the API interface by opening [http://localhost:8080/apis/swagger/index.html](http://localhost:8080/apis/swagger/index.html) in your browser.

![micro-cluster-swagger](https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/examples/micro-rpc-gw-pb-shopgw-swagger.png)

<br>
