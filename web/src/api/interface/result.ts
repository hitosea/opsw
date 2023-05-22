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
