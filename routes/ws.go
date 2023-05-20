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
	wsMsg      = vars.WsMsgStruct{}
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
	replyByte, _ := json.Marshal(map[string]any{
		"type": wsOnline,
		"data": map[string]any{
			"type": client.Type,
			"uid":  client.Uid,
			"rid":  client.Rid,
		},
	})
	_ = conn.WriteMessage(websocket.TextMessage, replyByte)
	// 循环读取客户端发送的消息
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, msg, err := conn.ReadMessage()
		if err != nil {
			app.wsOfflineClients(client.Rid)
			break
		}
		err = json.Unmarshal(msg, &wsMsg)
		if err != nil {
			continue
		}
		if wsMsg.Data == nil {
			wsMsg.Data = make(map[string]any)
		}
		replyByte = nil
		if wsMsg.Type == wsHeartbeat {
			// 心跳消息
			replyByte, _ = json.Marshal(map[string]any{
				"type": wsHeartbeat,
			})
		} else if client.Type == "user" {
			// 用户消息
			if wsMsg.Type == wsSendMsg {
				// 消息发送
				toType, _ := wsMsg.Data.(map[string]any)["to_type"].(string) // 客户端类型
				toUid, _ := wsMsg.Data.(map[string]any)["to_uid"].(float64)  // 发送给谁
				msgData, _ := wsMsg.Data.(map[string]any)["msg_data"].(any)  // 消息内容
				if toUid == 0 || msgData == nil {
					continue
				}
				if toType == "" {
					toType = "user"
				}
				sendByte, _ := json.Marshal(map[string]any{
					"type": wsSendMsg,
					"data": msgData,
				})
				for _, v := range vars.WsClients {
					if v.Type == toType && v.Uid == int32(toUid) {
						_ = v.Conn.WriteMessage(websocket.TextMessage, sendByte)
					}
				}
			} else if wsMsg.Type == wsOnlineClient {
				// 在线客户端
				var list []map[string]any
				for _, v := range vars.WsClients {
					list = append(list, map[string]any{
						"type": v.Type,
						"uid":  v.Uid,
						"rid":  v.Rid,
					})
				}
				replyByte, _ = json.Marshal(map[string]any{
					"type": wsOnlineClient,
					"data": map[string]any{
						"count": len(list),
						"list":  list,
					},
				})
			}
		} else if client.Type == "server" {
			// 服务器消息
		}
		if replyByte != nil {
			_ = conn.WriteMessage(websocket.TextMessage, replyByte)
		}
	}
}

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

func (app *AppStruct) wsNotifyClients(rid string, msgData vars.WsMsgStruct) {
	sendByte, err := json.Marshal(msgData)
	if err != nil {
		return
	}
	for _, client := range vars.WsClients {
		if client.Rid != rid {
			_ = client.Conn.WriteMessage(websocket.TextMessage, sendByte)
		}
	}
}
