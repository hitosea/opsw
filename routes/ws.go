package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"opsw/utils"
	"opsw/vars"
	"sync"
)

var (
	wsMutex    = sync.Mutex{}
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	wsHeartbeat    = 0 // 心跳
	wsOnline       = 1 // 连接
	wsOffline      = 2 // 断开
	wsSendMsg      = 3 // 消息发送
	wsOnlineClient = 4 // 获取在线客户端
)

// NoAuthWs 启动 websocket
func (app *AppStruct) NoAuthWs() {
	conn, err := wsUpgrader.Upgrade(app.Context.Writer, app.Context.Request, nil)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "连接失败", gin.H{"error": err.Error()})
		return
	}
	client := vars.WsClientStruct{
		Conn: conn,
	}
	if app.UserInfo.Id > 0 {
		client.Type = "user" // 用户
		client.Uid = app.UserInfo.Id
		client.Rid = fmt.Sprintf("u-%d-%s", client.Uid, utils.GenerateString(6))
	} else {
		client.Type = "server" // 服务器
		client.Uid = app.ServerInfo.Id
		client.Rid = fmt.Sprintf("s-%d-%s", client.Uid, utils.GenerateString(6))
	}
	// 完成时关闭连接释放资源
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-app.Context.Done()
		// 客户端离线
		app.wsOfflineClients(client.Rid)
	}()
	// 添加客户端（上线）
	app.wsOnlineClients(client)
	sendMsg, _ := json.Marshal(map[string]any{
		"type": wsOnline,
		"data": map[string]any{
			"type": client.Type,
			"uid":  client.Uid,
			"rid":  client.Rid,
			"own":  1,
		},
	})
	_ = conn.WriteMessage(websocket.TextMessage, sendMsg)
	// 循环读取客户端发送的消息
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, message, err := conn.ReadMessage()
		if err != nil {
			app.wsOfflineClients(client.Rid)
			break
		}
		var msg vars.WsMsgStruct
		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		}
		if msg.Data == nil {
			msg.Data = make(map[string]any)
		}
		if msg.Type == wsHeartbeat {
			// 心跳消息
			sendMsg, _ = json.Marshal(map[string]any{
				"type": wsHeartbeat,
			})
			_ = conn.WriteMessage(websocket.TextMessage, sendMsg)
			continue
		}
		if client.Type == "user" {
			// 用户消息
			app.wsHandleUserMsg(conn, msg)
		} else if client.Type == "server" {
			// 服务器消息
			app.wsHandleServerMsg(conn, msg)
		}
	}
}

// 处理用户消息
func (app *AppStruct) wsHandleUserMsg(conn *websocket.Conn, msg vars.WsMsgStruct) {
	var replyMsg []byte
	if msg.Type == wsSendMsg {
		// 消息发送
		toType, _ := msg.Data.(map[string]any)["to_type"].(string) // 客户端类型
		toUid, _ := msg.Data.(map[string]any)["to_uid"].(float64)  // 发送给谁
		msgData, _ := msg.Data.(map[string]any)["msg_data"].(any)  // 消息内容
		if toUid == 0 || msgData == nil {
			return
		}
		if toType == "" {
			toType = "user"
		}
		sendMsg, _ := json.Marshal(map[string]any{
			"type": wsSendMsg,
			"data": msgData,
		})
		for _, v := range vars.WsClients {
			if v.Type == toType && v.Uid == int32(toUid) {
				_ = v.Conn.WriteMessage(websocket.TextMessage, sendMsg)
			}
		}
	} else if msg.Type == wsOnlineClient {
		// 在线客户端
		var list []map[string]any
		for _, v := range vars.WsClients {
			list = append(list, map[string]any{
				"type": v.Type,
				"uid":  v.Uid,
				"rid":  v.Rid,
			})
		}
		replyMsg, _ = json.Marshal(map[string]any{
			"type": wsOnlineClient,
			"data": map[string]any{
				"count": len(list),
				"list":  list,
			},
		})
	}
	if replyMsg != nil {
		_ = conn.WriteMessage(websocket.TextMessage, replyMsg)
	}
}

// 处理服务器消息
func (app *AppStruct) wsHandleServerMsg(conn *websocket.Conn, msg vars.WsMsgStruct) {

}

// 客户端上线
func (app *AppStruct) wsOnlineClients(client vars.WsClientStruct) {
	for _, v := range vars.WsClients {
		if v.Rid == client.Rid {
			return
		}
	}
	wsMutex.Lock()
	vars.WsClients = append(vars.WsClients, client)
	wsMutex.Unlock()
	app.wsNotifyClients(client.Rid, vars.WsMsgStruct{
		Type: wsOnline,
		Data: map[string]any{
			"type": client.Type,
			"uid":  client.Uid,
			"rid":  client.Rid,
		},
	})
}

// 客户端离线
func (app *AppStruct) wsOfflineClients(rid string) {
	for k, client := range vars.WsClients {
		if client.Rid == rid {
			wsMutex.Lock()
			vars.WsClients = append(vars.WsClients[:k], vars.WsClients[k+1:]...)
			_ = client.Conn.Close()
			wsMutex.Unlock()
			//
			app.wsNotifyClients(rid, vars.WsMsgStruct{
				Type: wsOffline,
				Data: map[string]any{
					"type": client.Type,
					"uid":  client.Uid,
					"rid":  client.Rid,
				},
			})
			break
		}
	}
}

// 通知客户端
func (app *AppStruct) wsNotifyClients(rid string, msgData vars.WsMsgStruct) {
	sendMsg, err := json.Marshal(msgData)
	if err != nil {
		return
	}
	for _, client := range vars.WsClients {
		if client.Rid != rid {
			_ = client.Conn.WriteMessage(websocket.TextMessage, sendMsg)
		}
	}
}
