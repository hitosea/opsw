package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"opsw/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NoAuthWs 启动 websocket
func (app *AppStruct) NoAuthWs() {
	ws, err := upgrader.Upgrade(app.Context.Writer, app.Context.Request, nil)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "连接失败", gin.H{"error": err.Error()})
		return
	}
	// 完成时关闭连接释放资源
	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(ws)
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-app.Context.Done()
		// 这里也可以做用户在线/下线功能
		fmt.Println("ws lost connection")
	}()
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error") // 离线
			fmt.Println(err)
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
