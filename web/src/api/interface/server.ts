import {DatabaseBase, Page, PageReq} from "./base";

export namespace Server {
    export interface Log {
        log: string
    }

    export interface LogReq {
        id?: number
        ip?: string
    }

    export interface Base extends DatabaseBase {
        ip: string
        username: string
        password: string
        port: string
        remark: string
        state: string
        token: string
        systems: string
    }

    export interface CreateReq {
        ip: string
        username?: string
        password?: string
        port?: string
        remark?: string
    }

    export interface OperationReq extends LogReq {
        operation: string
    }

    export interface Item extends Base {
        server_id: number
        user_id: number
        owner_id: number

        hostname: string
        platform: string
        platform_version: string
        version: string

        upgrade: string
    }

    export interface List extends Page {
        data: Server.Item[]
    }

    export interface OneReq extends LogReq {
    }

    export interface ListReq extends PageReq {
        key?: string
    }
}
