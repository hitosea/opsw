import {User} from "../../api/interface/user";

export interface WsModel {
    ws: WebSocket,
    msg: WsMsgModel,
    uid: number,
    rid: string,
    timeout: any,
    random: string,
    openNum: number,
    listener: object,
}

export interface WsMsgModel {
    action?: number     // 消息类型：1、上线；2、下线；3、消息
    data?: any          // 消息内容

    type?: string       // 客户端类型：user、server
    cid?: number        // 客户端ID：用户ID、服务器ID
    rid?: string        // 客户端随机ID
}

export interface UserState {
    info: User.Info
}

export interface WsState {
    ws: WebSocket,
    msg: WsMsgModel,
    uid: number,
    rid: string,
    timeout: any,
    random: string,
    openNum: number,
    listener: object,
    watch: boolean
}
