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
        current_info: ItemInfoCurrent

        upgrade: string
    }


    export interface ItemInfoCurrent {
        uptime: number
        time_since_uptime: string

        procs: number

        load1: number
        load5: number
        load15: number
        load_usage_percent: number

        cpu_percent: number[]
        cpu_used_percent: number
        cpu_used: number
        cpu_total: number

        memory_total: number
        memory_available: number
        memory_used: number
        memory_used_percent: number

        io_read_bytes: number
        io_write_bytes: number
        io_count: number
        io_read_time: number
        io_write_time: number

        disk_data: ItemInfoDiskInfo[]

        net_bytes_sent: number
        net_bytes_recv: number

        shot_time: string
    }

    export interface ItemInfoDiskInfo {
        path: string
        type: string
        device: string
        total: number
        free: number
        used: number
        used_percent: number

        inodes_total: number
        inodes_used: number
        inodes_free: number
        inodes_used_percent: number
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
