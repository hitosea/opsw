export namespace User {
    export interface Info {
        id?: number
        email?: string
        name?: string
        token?: string
        avatar?: string
        created_at?: string
        updated_at?: string
    }

    export interface LoginReq {
        email: string,
        password: string,
    }

    export interface RegReq {
        email: string,
        password: string,
        password2?: string,
    }

    export interface Options {
        label: string,
        value: number,
        disabled: boolean
    }

    export interface ShareServerReq {
        server_id: number,
        user_ids: number[],
    }
}
