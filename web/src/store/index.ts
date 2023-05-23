import {computed, ref} from 'vue'
import {createPinia, defineStore} from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import result from "../utils/result";
import {GlobalState} from "./interface";
import piniaPersistConfig from "./config/pinia-persist";
import utils from "../utils/utils";

const nowTime = ref(0)
export function uptime(up: string|number, emptyTip = '-') {
    if (nowTime.value === 0) {
        nowTime.value = utils.Time()
        setInterval(() => {
            nowTime.value = utils.Time()
        }, 1000)
    }
    //
    if (!up) {
        return emptyTip
    }
    if (typeof up === 'string') {
        up = utils.Time(up)
    }
    let time: any = nowTime.value - up
    if (time === 0) {
        return `刚刚`
    }
    let day: any = Math.floor(time / 86400)
    let hour: any = Math.floor((time % 86400) / 3600)
    let minute: any = Math.floor((time % 3600) / 60)
    let second: any = time % 60
    if (day < 10) day = `0${day}`
    if (hour < 10) hour = `0${hour}`
    if (minute < 10) minute = `0${minute}`
    if (second < 10) second = `0${second}`
    if (day > 0) {
        return `${day}天${hour}时${minute}分`
    }
    if (hour > 0) {
        return `${hour}时${minute}分${second}秒`
    }
    if (minute > 0) {
        return `${minute}分${second}秒`
    }
    return `${second}秒`
}

export const GlobalStore = defineStore({
    id: 'GlobalState',
    state: (): GlobalState => ({
        isLoading: 0,
        themeName: '',
        timer: {},
    }),
    actions: {
        setLoading() {
            this.isLoading++;
        },
        cancelLoading() {
            this.isLoading--;
        },
        setThemeName(name: string) {
            this.themeName = name
        },
        async init() {
            this.isLoading = 0
            if (['light', 'dark'].indexOf(this.themeName) === -1) {
                this.themeName = useOsTheme().value
            }
        },
        appSetup() {
            return {
                resultCode: result.code(),
                resultMsg: result.msg(),
                themeName: this.themeName,
                theme: computed(() => {
                    return this.themeName === 'dark' ? darkTheme : null
                }),
            }
        },
        dialogProvider() {
            const {dialog} = createDiscreteApi(['dialog'], {
                configProviderProps: computed<ConfigProviderProps>(() => ({
                    theme: this.themeName === 'dark' ? darkTheme : null
                })),
            })
            return dialog
        },
        timeout(ms: number, key: string, ...name) {
            return new Promise(resolve => {
                key = `${key}-${name.join('-')}`
                this.timer[key] && clearTimeout(this.timer[key])
                if (typeof this.timer[key] !== "undefined") {
                    clearTimeout(this.timer[key])
                    delete this.timer[key]
                }
                this.timer[key] = setTimeout(resolve, ms)
            })
        }
    },
    persist: piniaPersistConfig('GlobalState'),
});

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
export default pinia;
