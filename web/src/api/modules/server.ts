import http from "../index";
import {Server} from "../interface/server";
import utils from "../../utils/utils";

export const getServerLog = (params: Server.LogReq) => {
    return http.get<Server.Log>('server/log', params)
}

export const createServer = (params: Server.CreateReq) => {
    return http.post<Server.Base>('server/create', params)
}

export const operationServer = (params: Server.OperationReq) => {
    return http.get('server/operation', params)
}

export const getServerList = (params: Server.ListReq) => {
    return http.get<Server.List>("server/list", params).then((res) => {
        res.data.data.forEach((item) => {
            item.current_info = utils.jsonParse(item.current_info)
        })
        return res
    })
}

export const getServerOne = (params: Server.OneReq) => {
    return http.get<Server.Item>("server/one", params).then((res) => {
        res.data.current_info = utils.jsonParse(res.data.current_info)
        return res
    })
}
