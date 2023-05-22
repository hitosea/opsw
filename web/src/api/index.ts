import axios, {AxiosInstance, AxiosRequestConfig} from 'axios'
import utils from "../utils/utils";
import {dialogProvider} from "../store";
import {ResultData} from "./interface/result";

const config = {
    baseURL: '/api', // 所有的请求地址前缀部分
    timeout: 60000, // 请求超时时间毫秒
    withCredentials: true, // 异步请求携带cookie
    headers: {
        // 设置后端需要的传参类型
        'Content-Type': 'application/x-www-form-urlencoded',
    },
}

class RequestHttp {
    // 定义成员变量并指定类型
    service: AxiosInstance;

    public constructor(config: AxiosRequestConfig) {
        // 实例化axios
        this.service = axios.create(config);

        /**
         * 请求拦截器
         * 客户端发送请求 -> [请求拦截器] -> 服务器
         * token校验(JWT) : 接受服务器返回的token,存储到vuex/pinia/本地储存当中
         */
        this.service.interceptors.request.use(
            function (config) {
                return config
            },
            function (error) {
                // 对请求错误做些什么
                // console.log(error)
                return Promise.reject(error)
            }
        )

        /**
         * 响应拦截器
         * 服务器换返回信息 -> [拦截统一处理] -> 客户端JS获取到信息
         */
        this.service.interceptors.response.use(
            function (response) {
                // console.log(response)
                // 2xx 范围内的状态码都会触发该函数。
                // 对响应数据做点什么
                // dataAxios 是 axios 返回数据中的 data
                const dataAxios = response.data
                //
                if (!utils.isJson(dataAxios)) {
                    return Promise.reject({code: 500, msg: "返回数据格式错误", data: dataAxios})
                }
                if (dataAxios.code !== 200) {
                    if (dataAxios.code === 301 || dataAxios.code === 401) {
                        const params = {
                            result_code: dataAxios.code,
                            result_msg: encodeURIComponent(dataAxios.msg),
                        }
                        window.location.href = utils.urlAddParams(window.location.href, params)
                    }
                    return Promise.reject(dataAxios)
                }
                return dataAxios
            },
            function (error) {
                // 超出 2xx 范围的状态码都会触发该函数。
                // 对响应错误做点什么
                // console.log(error)
                return Promise.reject({code: 500, msg: "请求失败", data: error})
            }
        )
    }

    // 自定义方法封装（常用请求）
    get<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.get(url, {params});
    }

    post<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.post(url, params);
    }

    put<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.put(url, params);
    }

    delete<T>(url: string, params?: object): Promise<ResultData<T>> {
        return this.service.delete(url, {params});
    }

    // 自定义方法封装（提示框）
    dialog({code, msg, data}, dialogOptions = {}) {
        let title = '温馨提示'
        let content = msg
        if (code !== 200) {
            title = '错误提示'
            if (utils.isJson(data) && (data.err || data.error)) {
                title = msg
                content = data.err || data.error
            }
        }
        let options = {
            title,
            content,
            positiveText: '确定',
        }
        if (utils.isJson(data) && utils.isJson(data.dialog)) {
            options = Object.assign(options, data.dialog)
        }
        if (utils.isJson(dialogOptions)) {
            options = Object.assign(options, dialogOptions)
        }
        if (code === 200) {
            return dialogProvider().success(options)
        } else {
            return dialogProvider().error(options)
        }
    }
}

export default new RequestHttp(config);
