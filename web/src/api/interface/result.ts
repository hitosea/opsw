export interface Result {
    code: number;
    msg: string
}

export interface ResultData<T = any> extends Result {
    data?: T;
}
