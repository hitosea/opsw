package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"opsw/utils"
	"opsw/vars"
	"sync"
)

var (
	clients   []vars.WsClientStruct
	clientMsg = vars.WsMsgStruct{}
	mutex     = sync.Mutex{}
)

const (
	heartbeat = 0 // 心跳
	online    = 1 // 连接
	offline   = 2 // 断开
	sendMsg   = 3 // 消息发送
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// AuthWs 启动 websocket
func (app *AppStruct) AuthWs() {
	conn, err := upgrader.Upgrade(app.Context.Writer, app.Context.Request, nil)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "连接失败", gin.H{"error": err.Error()})
		return
	}
	randId := utils.GenerateString(16)
	// 完成时关闭连接释放资源
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-app.Context.Done()
		// 这里也可以做用户在线/下线功能
		app.removeClients(randId)
	}()
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, message, err := conn.ReadMessage()
		if err != nil {
			app.removeClients(randId)
			break
		}
		_ = json.Unmarshal(message, &clientMsg)
		if clientMsg.Data == nil {
			continue
		}
		if clientMsg.State == heartbeat {
			// 心跳
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"state":0,"data":"heartbeat ok"}`))
		} else if clientMsg.State == online {
			// 连接
			app.addClients(randId, conn)
		} else if clientMsg.State == sendMsg {
			// 消息发送
			userId := clientMsg.Data.(map[string]any)["user_id"].(int32)    // 发送给谁
			msgData := clientMsg.Data.(map[string]any)["msg_data"].(string) // 消息内容
			msgByte, _ := json.Marshal(msgData)
			for _, v := range clients {
				if v.UserId == userId {
					_ = v.Conn.WriteMessage(websocket.TextMessage, msgByte)
				}
			}
		}
	}
}

func (app *AppStruct) addClients(randId string, conn *websocket.Conn) {
	for _, v := range clients {
		if v.RandId == randId {
			return
		}
	}
	mutex.Lock()
	clients = append(clients, vars.WsClientStruct{
		Conn:   conn,
		UserId: app.UserInfo.Id,
		RandId: randId,
	})
	mutex.Unlock()
	app.notifyClients(randId, vars.WsMsgStruct{
		State: online,
		Data: map[string]any{
			"user_id": app.UserInfo.Id,
		},
	})
}

func (app *AppStruct) removeClients(randId string) {
	for k, v := range clients {
		if v.RandId == randId {
			mutex.Lock()
			clients = append(clients[:k], clients[k+1:]...)
			_ = v.Conn.Close()
			mutex.Unlock()
			app.notifyClients(randId, vars.WsMsgStruct{
				State: offline,
				Data: map[string]any{
					"user_id": app.UserInfo.Id,
				},
			})
			break
		}
	}
}

func (app *AppStruct) notifyClients(randId string, msgData vars.WsMsgStruct) {
	msgByte, err := json.Marshal(msgData)
	if err != nil {
		return
	}
	for _, v := range clients {
		if v.RandId != randId {
			_ = v.Conn.WriteMessage(websocket.TextMessage, msgByte)
		}
	}
}
