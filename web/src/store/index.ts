import {computed, watch, ref} from 'vue'
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import utils from "../utils/utils";
import result from "../utils/result";
import local from "../utils/local";
import {User} from "../api/interface/user";
import {getUserInfo} from "../api/modules/user";
import {WsModel, WsMsgModel} from "./interface";
import {CONST} from "./constant";


const wsRef = ref<WsModel>({
    uid: 0,
    rid: null,
    ws: null,
    msg: {},
    timeout: null,
    random: "",
    openNum: 0,
    listener: {},
})
const userInfoRef = ref<User.Info>({})
const themeNameRef = ref('light')
const themeRef = computed(() => {
    const {value} = themeNameRef
    return value === 'dark' ? darkTheme : null
})
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
    theme: themeRef.value
}))


watch(themeNameRef, name => {
    local.save("themeName", name)
})

watch(userInfoRef, info => {
    if (wsRef.value.uid !== info.id) {
        wsRef.value.uid = info.id
        wsConnection()
    }
})

export function useThemeName() {
    return themeNameRef
}

export function useUserInfo() {
    return userInfoRef
}

export function useWs() {
    return wsRef
}

export function loadUserInfo() {
    return new Promise((resolve, reject) => {
        getUserInfo()
            .then(({data}) => {
                if (utils.isEmpty(data.name)) {
                    data.name = data.email
                }
                userInfoRef.value = data
                resolve(data)
            })
            .catch(err => {
                reject(err)
            })
    })
}

export function wsConnection() {
    clearTimeout(wsRef.value.timeout);
    if (wsRef.value.ws) {
        wsRef.value.ws.close();
        wsRef.value.ws = null;
    }
    if (userInfoRef.value.id === 0) {
        return;
    }
    //
    let url = window.location.origin + "/ws";
    url = url.replace("https://", "wss://");
    url = url.replace("http://", "ws://");
    url += `?token=${userInfoRef.value.token}`;
    //
    const wgLog = true;
    const random = utils.randomString(16);
    wsRef.value.random = random;
    //
    wsRef.value.ws = new WebSocket(url);
    wsRef.value.ws.onopen = async (e) => {
        wgLog && console.log("[WS] Open", e, utils.formatDate())
        wsRef.value.openNum++;
    };
    wsRef.value.ws.onclose = async (e) => {
        wgLog && console.log("[WS] Close", e, utils.formatDate())
        wsRef.value.ws = null;
        //
        clearTimeout(wsRef.value.timeout);
        wsRef.value.timeout = setTimeout(() => {
            random === wsRef.value.random && wsConnection();
        }, 3000);
    };
    wsRef.value.ws.onerror = async (e) => {
        wgLog && console.log("[WS] Error", e, utils.formatDate())
        wsRef.value.ws = null;
        //
        clearTimeout(wsRef.value.timeout);
        wsRef.value.timeout = setTimeout(() => {
            random === wsRef.value.random && wsConnection();
        }, 3000);
    };
    wsRef.value.ws.onmessage = async (e) => {
        wgLog && console.log("[WS] Message", e);
        const wsMsg: WsMsgModel = utils.jsonParse(e.data);
        wsRef.value.msg = wsMsg;
        //
        if (wsMsg.action === CONST.WsOnline) {
            // 连接成功
            if (wsMsg.data.own === 1) {
                wsRef.value.rid = wsMsg.data.rid
            }
        }
        Object.values(wsRef.value.listener).forEach(call => {
            if (typeof call === "function") {
                try {
                    call(wsMsg);
                } catch (err) {
                    wgLog && console.log("[WS] Callerr", err);
                }
            }
        });
    }
}

export function wsMsgListener(name, callback) {
    if (typeof callback === "function") {
        wsRef.value.listener[name] = callback;
    } else {
        wsRef.value.listener[name] && delete wsRef.value.listener[name];
    }
}

export function wsClose() {
    if (wsRef.value.ws) {
        wsRef.value.ws.close();
        wsRef.value.ws = null;
    }
}

export function dialogProvider() {
    const {dialog} = createDiscreteApi(['dialog'], {
        configProviderProps: configProviderPropsRef,
    })
    return dialog
}

export function siteSetup() {
    return {
        resultCode: result.code(),
        resultMsg: result.msg(),
        themeName: themeNameRef,
        theme: themeRef,
    }
}

export function init() {
    return new Promise<void>(async (resolve) => {
        themeNameRef.value = <string>await local.string("themeName")
        if (['light', 'dark'].indexOf(themeNameRef.value) === -1) {
            themeNameRef.value = useOsTheme().value
        }
        resolve()
    })
}
