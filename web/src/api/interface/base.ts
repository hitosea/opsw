export interface Result {
    code: number;
    msg: string
}

export interface ResultData<T = any> extends Result {
    data?: T;
}

export interface DatabaseBase {
    id: number;
    created_at: number
    updated_at: number
}

export interface PageReq {
    page: number;
    page_size?: number
}

export interface Page {
    page: number
    page_size: number

    next_page: number
    prev_page: number
    page_count: number

    total: number
}
