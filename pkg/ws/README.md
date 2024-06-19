## WebSocket

`ws` is based on the [github.com/gorilla/websocket](https://github.com/gorilla/websocket) library.

<br>

### Example of use

#### 1. default setting

**Server side code example:**

```go
package main

import (
    "context"
    "log"
    "github.com/zhufuyi/sponge/pkg/ws"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
	
    r.GET("/ws", func(c *gin.Context) {
        s := ws.NewServer(c.Writer, c.Request, loopReceiveMessage) // default setting
        err := s.Run(context.Background())
        if err != nil {
            log.Println("webSocket server error:", err)
        }
    })
	
    err := r.Run(":8080")
    if err != nil {
        panic(err)
    }
}

func loopReceiveMessage(ctx context.Context, conn *ws.Conn) {
    for {
        messageType, message, err := conn.ReadMessage()
        // handle message
        log.Println(messageType, message, err)
    }
}
```

<br>

**Client side code example:**

```go
package main

import (
    "strconv"
    "log"
    "time"
    "github.com/zhufuyi/sponge/pkg/ws"
    "github.com/gorilla/websocket"
)

var wsURL = "ws://localhost:8080/ws"

func main() {
    c, err := ws.NewClient(wsURL)
    if err != nil {
        log.Println("connect error:", err)
        return
    }
    defer c.Close()

    go func() {
        for {
            _, message, err := c.GetConn().ReadMessage()
            if err != nil {
                log.Println("client read error:", err)
                return
            }
            log.Printf("client received: %s", message)
        }
    }()
    
    for i := 0; i < 5; i++ {
        data := "Hello, World " + strconv.Itoa(i)
        err = c.GetConn().WriteMessage(websocket.TextMessage, []byte(data))
        if err != nil {
            log.Println("write error:", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

<br>

#### 2. custom setting

**Server side custom setting**, options such as `ws.Upgrader` `ws.WithMaxMessageWaitPeriod` can be set.

```go
package main

import (
    "context"
    "log"
    "time"
    "http"
    "github.com/zhufuyi/sponge/pkg/ws"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

func main() {
    r := gin.Default()
    ug := &websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }    
    r.GET("/ws", func(c *gin.Context) {
        s := ws.NewServer(c.Writer, c.Request, loopReceiveMessage,
            ws.WithUpgrader(ug),
            ws.WithMaxMessageWaitPeriod(time.Minute),
        )
        err := s.Run(context.Background())
        if err != nil {
            log.Println("webSocket server error:", err)
        }
    })
    
    err := r.Run(":8080")
    if err != nil {
        panic(err)
    }
}

func loopReceiveMessage(ctx context.Context, conn *ws.Conn) {
    for {
        messageType, message, err := conn.ReadMessage()
        // handle message
        log.Println(messageType, message, err)
    }
}
```

<br>

**Client side custom setting**, options such as `ws.Dialer` `ws.WithPingInterval` can be set.

```go
package main

import (
    "strconv"
    "log"
    "time"
    "github.com/zhufuyi/sponge/pkg/ws"
    "github.com/gorilla/websocket"
)

var wsURL = "ws://localhost:8080/ws"

func main() {
    c, err := NewClient(wsURL,
        WithDialer(websocket.DefaultDialer),
        WithPing(time.Second*10),
    )
    if err != nil {
        log.Println("connect error:", err)
        return
    }
    defer c.Close()

    go func() {
        for {
            _, message, err := c.GetConn().ReadMessage()
            if err != nil {
                log.Println("client read error:", err)
                return
            }
            log.Printf("client received: %s", message)
        }
    }()

    for i := 5; i < 10; i++ {
        data := "Hello, World " + strconv.Itoa(i)
        err = c.GetConn().WriteMessage(websocket.TextMessage, []byte(data))
        if err != nil {
            log.Println("write error:", err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    <-time.After(time.Minute)
}
```
