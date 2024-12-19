## swagger

gin swagger library.

<br>

## Example of use

```go
package main

import (
    "net/http"

    "github.com/go-dev-frame/sponge/pkg/gin/validator"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
)

func main() {
	r := gin.Default()
	binding.Validator = validator.Init()
	
	r.Run(":8080")
}

type createUserRequest struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Age   int    `json:"age" form:"age" binding:"gte=0,lte=120"`
	Email string `json:"email" form:"email" binding:"email"`
}

func CreateUser(c *gin.Context) {
	form := &createUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

type getUserRequest struct {
	Page int    `json:"page" form:"page" binding:"gte=0"`
	Limit int    `json:"limit" form:"limit" binding:"gte=1"`
	Sort string `json:"sort" form:"sort" binding:"-"`
}

func GetUsers(c *gin.Context) {
	form := &getUserRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	users, err := getUsers(form)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
        return
    }

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func getUsers(req *getUserRequest) ([]User,error){}
```


