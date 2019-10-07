package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   []*websocket.Conn
	nwsPeriod = 10 * time.Second
)

type nodeMessage struct {
	Hostname    string  `json:"hostname"`
	Platform    string  `json:"platform"`
	Kernel      string  `json:"kernel"`
	Uptime      uint64  `json:"uptime"`
	Load1       float64 `json:"load1"`
	Load5       float64 `json:"load5"`
	Load15      float64 `json:"load15"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
}

func index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func nodeWebsocket(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		time.Sleep(nwsPeriod)
		if len(clients) == 0 {
			continue
		}
		if err := ws.WriteMessage(websocket.TextMessage, []byte("check")); err != nil {
			return
		}
		var message nodeMessage
		if err := ws.ReadJSON(&message); err != nil {
			return
		}
		for _, client := range clients {
			if err := client.WriteJSON(message); err != nil {
				clients = pop(clients, client)
			}
		}
	}
}

func clientWebsocket(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	clients = append(clients, ws)
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/", index)
	r.GET("/nws", nodeWebsocket)
	r.GET("/cws", clientWebsocket)
	r.Run()
}

func pop(clients []*websocket.Conn, t *websocket.Conn) []*websocket.Conn {
	for idx, client := range clients {
		if client == t {
			return append(clients[0:idx], clients[idx+1:]...)
		}
	}
	return []*websocket.Conn{}
}
