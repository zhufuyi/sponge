package ws

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/go-dev-frame/sponge/pkg/utils"
)

func TestWebSocketServerDefault(t *testing.T) {
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		s := NewServer(c.Writer, c.Request, loopReceiveMessage)
		err := s.Run(context.Background())
		if err != nil {
			log.Println("WebSocket server error:", err)
		}
	})

	addr, _ := utils.GetLocalHTTPAddrPairs()
	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	url := "ws://localhost" + addr + "/ws"
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			runWsClientDefault(url)
		}()
	}
	wg.Wait()

	time.Sleep(100 * time.Millisecond)
}

func TestWebSocketServerCustom(t *testing.T) {
	r := gin.Default()
	ug := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	l, _ := zap.NewProduction()
	r.GET("/ws", func(c *gin.Context) {
		s := NewServer(c.Writer, c.Request, loopReceiveMessage,
			WithUpgrader(ug),
			WithResponseHeader(http.Header{"Authorization": []string{"Bearer 123"}}),
			WithNoClientPingTimeout(time.Second*3),
			WithServerLogger(l),
		)
		err := s.Run(context.Background())
		if err != nil {
			log.Println("WebSocket server error:", err)
		}
	})

	addr, _ := utils.GetLocalHTTPAddrPairs()
	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	url := "ws://localhost" + addr + "/ws"
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			runWsClientCustom(url)
		}()
	}
	wg.Wait()

	time.Sleep(100 * time.Millisecond)
}

func loopReceiveMessage(ctx context.Context, conn *Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if IsClientClose(err) {
				log.Println("Client closed connection")
				return
			}
			log.Println("Read message error:", err)
			return
		}
		log.Printf("server side received: %s", message)

		switch messageType {
		case websocket.TextMessage:
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Write message error:", err)
				continue
			}
		case websocket.PingMessage:
			err = conn.WriteMessage(websocket.PongMessage, []byte("pong"))
			if err != nil {
				log.Println("Write pong message error:", err)
				continue
			}
		default:
			log.Printf("Unknown message type: %d", messageType)
		}
	}
}

var wsURL = "ws://localhost:8080/ws"

func runWsClientDefault(url string) {
	if url == "" {
		url = wsURL
	}
	c, err := NewClient(url)
	if err != nil {
		log.Println("connect error:", err)
		return
	}
	defer c.Close()

	go clientLoopReadMessage(c.GetConn())

	for i := 0; i < 10; i++ {
		data := "Hello, World " + strconv.Itoa(i)
		err = c.GetConn().WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("write error:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func runWsClientCustom(url string) {
	if url == "" {
		url = wsURL
	}
	l, _ := zap.NewProduction()
	c, err := NewClient(url,
		WithDialer(websocket.DefaultDialer),
		WithRequestHeader(http.Header{"foo": []string{"bar"}}),
		WithPing(time.Second),
		WithClientLogger(l),
	)
	if err != nil {
		log.Println("connect error:", err)
		return
	}
	defer c.Close()

	go clientLoopReadMessage(c.GetConn())

	for i := 5; i < 10; i++ {
		data := "Hello, World " + strconv.Itoa(i)
		err = c.GetConn().WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Println("write error:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	<-time.After(time.Second * 2)
	_ = c.GetConn().WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
	<-time.After(time.Second * 2)
}

func clientLoopReadMessage(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			return
		}
		log.Printf("client side received: %s", message)
	}
}
