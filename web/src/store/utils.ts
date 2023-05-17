import Cookies from "js-cookie";
import * as localforage from "localforage";

localforage.config({name: 'web', storeName: 'common'});

const utils = {
    /**
     * 是否数组
     * @param obj
     * @returns {boolean}
     */
    isArray(obj) {
        return typeof (obj) == "object" && Object.prototype.toString.call(obj).toLowerCase() == '[object array]' && typeof obj.length == "number";
    },

    /**
     * 是否数组对象
     * @param obj
     * @returns {boolean}
     */
    isJson(obj) {
        return typeof (obj) == "object" && Object.prototype.toString.call(obj).toLowerCase() == "[object object]" && typeof obj.length == "undefined";
    },

    /**
     * 获取对象值
     * @param obj
     * @param key
     * @returns {*}
     */
    getObject(obj, key) {
        const keys = key.replace(/,/g, "|").replace(/\./g, "|").split("|");
        while (keys.length > 0) {
            const k = keys.shift();
            if (utils.isArray(obj)) {
                obj = obj[utils.parseInt(k)] || "";
            } else if (utils.isJson(obj)) {
                obj = obj[k] || "";
            } else {
                break;
            }
        }
        return obj;
    },

    /**
     * 转成数字
     * @param param
     * @returns {number|number}
     */
    parseInt(param) {
        const num = parseInt(param);
        return isNaN(num) ? 0 : num;
    },

    /**
     * 是否在数组里
     * @param key
     * @param array
     * @param regular
     * @returns {boolean|*}
     */
    inArray(key, array, regular = false) {
        if (!utils.isArray(array)) {
            return false;
        }
        if (regular) {
            return !!array.find(item => {
                if (item && item.indexOf("*")) {
                    const rege = new RegExp("^" + item.replace(/[-\/\\^$+?.()|[\]{}]/g, '\\$&').replace(/\*/g, '.*') + "$", "g")
                    if (rege.test(key)) {
                        return true
                    }
                }
                return item == key
            });
        } else {
            return array.includes(key);
        }
    },

    /**
     * 克隆对象
     * @param myObj
     * @returns {*}
     */
    cloneJSON(myObj) {
        if (typeof (myObj) !== 'object') return myObj;
        if (myObj === null) return myObj;
        //
        return utils.jsonParse(utils.jsonStringify(myObj))
    },

    /**
     * 将一个 JSON 字符串转换为对象（已try）
     * @param str
     * @param defaultVal
     * @returns {*}
     */
    jsonParse(str, defaultVal = undefined) {
        if (str === null) {
            return defaultVal ? defaultVal : {};
        }
        if (typeof str === "object") {
            return str;
        }
        try {
            return JSON.parse(str.replace(/\n/g, "\\n").replace(/\r/g, "\\r"));
        } catch (e) {
            return defaultVal ? defaultVal : {};
        }
    },

    /**
     * 将 JavaScript 值转换为 JSON 字符串（已try）
     * @param json
     * @param defaultVal
     * @returns {string}
     */
    jsonStringify(json, defaultVal = undefined) {
        if (typeof json !== 'object') {
            return json;
        }
        try {
            return JSON.stringify(json);
        } catch (e) {
            return defaultVal ? defaultVal : "";
        }
    },

    /**
     * 字符串是否包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    strExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.indexOf(find) !== -1);
    },

    /**
     * 字符串是否左边包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    leftExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.substring(0, find.length) === find);
    },

    /**
     * 删除左边字符串
     * @param string
     * @param find
     * @param lower
     * @returns {string}
     */
    leftDelete(string, find, lower = false) {
        string += "";
        find += "";
        if (utils.leftExists(string, find, lower)) {
            string = string.substring(find.length)
        }
        return string ? string : '';
    },

    /**
     * 字符串是否右边包含
     * @param string
     * @param find
     * @param lower
     * @returns {boolean}
     */
    rightExists(string, find, lower = false) {
        string += "";
        find += "";
        if (lower !== true) {
            string = string.toLowerCase();
            find = find.toLowerCase();
        }
        return (string.substring(string.length - find.length) === find);
    },

    /**
     * 删除右边字符串
     * @param string
     * @param find
     * @param lower
     * @returns {string}
     */
    rightDelete(string, find, lower = false) {
        string += "";
        find += "";
        if (utils.rightExists(string, find, lower)) {
            string = string.substring(0, string.length - find.length)
        }
        return string ? string : '';
    },

    /**
     * 检测手机号码格式
     * @param str
     * @returns {boolean}
     */
    isPhone(str) {
        return /^1([3456789])\d{9}$/.test(str);
    },

    /**
     * 检测邮箱地址格式
     * @param email
     * @returns {boolean}
     */
    isEmail(email) {
        return /^([0-9a-zA-Z]([-.\w]*[0-9a-zA-Z])*@([0-9a-zA-Z][-\w]*\.)+[a-zA-Z]*)$/i.test(email);
    },

    /**
     * 指定键获取url参数
     * @param key
     * @returns {*}
     */
    urlParameter(key) {
        const params = utils.urlParameterAll();
        return typeof key === "undefined" ? params : params[key];
    },

    urlParameterAll() {
        let search = window.location.search || window.location.hash || "";
        const index = search.indexOf("?");
        if (index !== -1) {
            search = search.substring(index + 1);
        }
        const arr = search.split("&");
        const params = {};
        arr.forEach((item) => { // 遍历数组
            const index = item.indexOf("=");
            if (index === -1) {
                params[item] = "";
            } else {
                params[item.substring(0, index)] = item.substring(index + 1);
            }
        });
        return params;
    },

    /**
     * 删除地址中的参数
     * @param url
     * @param parameter
     * @returns {string|*}
     */
    removeURLParameter(url, parameter) {
        if (parameter instanceof Array) {
            parameter.forEach((key) => {
                url = utils.removeURLParameter(url, key)
            });
            return url;
        }
        const urlParts = url.split('?');
        if (urlParts.length >= 2) {
            //参数名前缀
            let prefix = encodeURIComponent(parameter) + '=';
            let pars = urlParts[1].split(/[&;]/g);

            //循环查找匹配参数
            for (let i = pars.length; i-- > 0;) {
                if (pars[i].lastIndexOf(prefix, 0) !== -1) {
                    //存在则删除
                    pars.splice(i, 1);
                }
            }

            return urlParts[0] + (pars.length > 0 ? '?' + pars.join('&') : '');
        }
        return url;
    },

    /**
     * 连接加上参数
     * @param url
     * @param params
     * @returns {*}
     */
    urlAddParams(url, params) {
        if (utils.isJson(params)) {
            if (url) {
                url = utils.removeURLParameter(url, Object.keys(params))
            }
            url += "";
            url += url.indexOf("?") === -1 ? '?' : '';
            for (let key in params) {
                if (!params.hasOwnProperty(key)) {
                    continue;
                }
                url += '&' + key + '=' + params[key];
            }
        } else if (params) {
            url += (url.indexOf("?") === -1 ? '?' : '&') + params;
        }
        if (!url) {
            return ""
        }
        return utils.rightDelete(url.replace("?&", "?"), '?');
    },

    /**
     * 获取返回码
     * @returns {*|number}
     */
    resultCode() {
        return utils.parseInt(utils.urlParameter("result_code") || window['result_code'])
    },

    /**
     * 获取返回消息
     * @returns {*|string|string}
     */
    resultMsg() {
        return decodeURIComponent(utils.urlParameter("result_msg") || window['result_msg'])
    },

    /**
     * =============================================================================
     * ********************************   cookie   *********************************
     * =============================================================================
     */

    /**
     * 获取cookie
     * @param name
     * @param defaultVal
     * @returns {string}
     * @constructor
     */
    GetCookie(name, defaultVal = "") {
        return decodeURIComponent(Cookies.get(name)) || defaultVal
    },

    /**
     * 设置cookie
     * @param name
     * @param value
     * @constructor
     */
    SetCookie(name, value) {
        Cookies.set(name, encodeURIComponent(value))
    },

    /**
     * 删除cookie
     * @param name
     * @constructor
     */
    RemoveCookie(name) {
        Cookies.remove(name)
    },

    /**
     * =============================================================================
     * *****************************   localForage   ******************************
     * =============================================================================
     */
    __IDBTimer: {},

    IDBSave(key, value, delay = 100) {
        if (typeof utils.__IDBTimer[key] !== "undefined") {
            clearTimeout(utils.__IDBTimer[key])
            delete utils.__IDBTimer[key]
        }
        utils.__IDBTimer[key] = setTimeout(async _ => {
            await localforage.setItem(key, value)
        }, delay)
    },

    IDBDel(key) {
        localforage.removeItem(key).then(_ => {
        })
    },

    IDBSet(key, value) {
        return localforage.setItem(key, value)
    },

    IDBRemove(key) {
        return localforage.removeItem(key)
    },

    IDBClear() {
        return localforage.clear()
    },

    IDBValue(key) {
        return localforage.getItem(key)
    },

    async IDBString(key, def = "") {
        const value = await utils.IDBValue(key)
        return typeof value === "string" || typeof value === "number" ? value : def;
    },

    async IDBInt(key, def = 0) {
        const value = await utils.IDBValue(key)
        return typeof value === "number" ? value : def;
    },

    async IDBBoolean(key, def = false) {
        const value = await utils.IDBValue(key)
        return typeof value === "boolean" ? value : def;
    },

    async IDBArray(key, def = []) {
        const value = await utils.IDBValue(key)
        return utils.isArray(value) ? value : def;
    },

    async IDBJson(key, def = {}) {
        const value = await utils.IDBValue(key)
        return utils.isJson(value) ? value : def;
    }
}

export default utils
