import http from "../index";
import {Server} from "../interface/server";

export const getServerLog = (params: Server.LogReq) => {
    return http.get<Server.Log>('server/log', params)
}

export const createServer = (params: Server.CreateReq) => {
    return http.post<Server.Base>('server/create', params)
}

export const operationServer = (params: Server.OperationReq) => {
    return http.get('server/operation', params)
}

export const getServerList = () => {
    return http.get<Server.List>("server/list")
}

export const getServerOne = (params: Server.OneReq) => {
    return http.get<Server.Item>("server/one", params)
}
